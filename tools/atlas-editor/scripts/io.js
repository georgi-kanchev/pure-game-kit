// Save
let savedFileHandle = null;

document.addEventListener('keydown', async (e) => {
    if (e.ctrlKey && e.key === 's') {
        e.preventDefault();
        if (savedFileHandle) {
            const writable = await savedFileHandle.createWritable();
            await writable.write(buildXml());
            await writable.close();
        } else {
            exportXml();
        }
    }
});

document.getElementById('save').addEventListener('click', exportXml);

async function exportXml() {
    const content = buildXml();
    if (window.showSaveFilePicker) {
        try {
            const handle = await window.showSaveFilePicker({
                suggestedName: 'atlas.xml',
                types: [{ description: 'XML', accept: { 'application/xml': ['.xml'] } }],
            });
            savedFileHandle = handle;
            const writable = await handle.createWritable();
            await writable.write(content);
            await writable.close();
            return;
        } catch (e) {
            if (e.name === 'AbortError') return;
        }
    }
    const blob = new Blob([content], { type: 'application/xml' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'atlas.xml';
    a.click();
    URL.revokeObjectURL(url);
}

function buildXml() {
    const lines = ['<?xml version="1.0" encoding="UTF-8"?>', `<data grid="${gridSize}">`];

    lines.push('  <crops>');
    crops.forEach((f, i) => {
        lines.push(`    <crop x="${f.x}" y="${f.y}" w="${f.w}" h="${f.h}"/>`);
    });
    lines.push('  </crops>');

    lines.push('  <groups>');
    groups.forEach((a, i) => {
        const indices = a.cropIndices.join(' ');
        lines.push(`    <group name=${xmlAttr(a.name)} cropIndexes="${indices}"/>`);
    });
    lines.push('  </groups>');

    lines.push('</data>');
    return lines.join('\n');
}

// Load
document.getElementById('load').addEventListener('click', () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.xml';
    input.addEventListener('change', e => {
        const file = e.target.files[0];
        if (!file) return;
        const reader = new FileReader();
        reader.onload = ev => importXml(ev.target.result);
        reader.readAsText(file);
    });
    input.click();
});

function importXml(text) {
    const doc = new DOMParser().parseFromString(text, 'application/xml');
    if (doc.querySelector('parsererror')) return;

    // load grid size
    const dataEl = doc.querySelector('data');
    if (dataEl) {
        const g = parseInt(dataEl.getAttribute('grid'));
        if (!isNaN(g) && g > 0) {
            gridSize = g;
            document.getElementById('gridSize').value = g;
            if (image) buildChecker(image.width, image.height);
        }
    }

    // load crops — generate hues
    const cropEls = [...doc.querySelectorAll('crops > crop')];
    crops.length = 0;
    let chue = Math.random() * 360;
    const seen = new Set();
    cropEls.forEach(el => {
        const x = parseFloat(el.getAttribute('x')) || 0;
        const y = parseFloat(el.getAttribute('y')) || 0;
        const w = parseFloat(el.getAttribute('w')) || 0;
        const h = parseFloat(el.getAttribute('h')) || 0;
        const key = `${x},${y},${w},${h}`;
        if (seen.has(key)) return; // skip exact duplicate crops
        seen.add(key);
        crops.push({ x, y, w, h, hue: chue });
        chue = (chue + 20 + Math.random() * 25) % 360;
    });
    lastHue = chue;

    // load groups — generate hues
    const groupEls = [...doc.querySelectorAll('groups > group')];
    groups.length = 0;
    selectedGroupIdx = -1;
    let ahue = Math.random() * 360;
    groupEls.forEach(el => {
        const indices = (el.getAttribute('cropIndexes') || '').trim().split(/\s+/).filter(Boolean).map(Number);
        groups.push({
            name: el.getAttribute('name') || '',
            hue: ahue,
            cropIndices: indices.filter(n => !isNaN(n)),
        });
        ahue = (ahue + 20 + Math.random() * 25) % 360;
    });
    if (groups.length) selectedGroupIdx = 0;

    rebuildGroupList();
    selectGroup(selectedGroupIdx);
    drawView();
}

function xmlAttr(v) {
    return `"${String(v).replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;')}"`;
}

// Open image
document.getElementById('openImage').addEventListener('click', () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = 'image/png';
    input.addEventListener('change', e => {
        const file = e.target.files[0];
        if (!file) return;
        const reader = new FileReader();
        reader.onload = ev => {
            const img = new Image();
            img.onload = () => {
                image = img;
                buildChecker(image.width, image.height);
                resetView();
            };
            img.src = ev.target.result;
        };
        reader.readAsDataURL(file);
    });
    input.click();
});

// Grid size
const gridInput = document.getElementById('gridSize');
gridInput.addEventListener('input', () => {
    gridSize = parseInt(gridInput.value) || 16;
    if (image) buildChecker(image.width, image.height);
    drawView();
});
