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
                suggestedName: 'animations.xml',
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
    a.download = 'animations.xml';
    a.click();
    URL.revokeObjectURL(url);
}

function buildXml() {
    const lines = ['<?xml version="1.0" encoding="UTF-8"?>', '<data>'];

    lines.push('  <frames>');
    frames.forEach((f, i) => {
        lines.push(`    <frame x="${f.x}" y="${f.y}" w="${f.w}" h="${f.h}"/>`);
    });
    lines.push('  </frames>');

    lines.push('  <animations>');
    animations.forEach((a, i) => {
        const indices = a.frameIndices.join(' ');
        lines.push(`    <animation name=${xmlAttr(a.name)} frames="${indices}"/>`);
    });
    lines.push('  </animations>');

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

    // load frames — generate hues
    const frameEls = [...doc.querySelectorAll('frames > frame')];
    frames.length = 0;
    let fhue = Math.random() * 360;
    const seen = new Set();
    frameEls.forEach(el => {
        const x = parseFloat(el.getAttribute('x')) || 0;
        const y = parseFloat(el.getAttribute('y')) || 0;
        const w = parseFloat(el.getAttribute('w')) || 0;
        const h = parseFloat(el.getAttribute('h')) || 0;
        const key = `${x},${y},${w},${h}`;
        if (seen.has(key)) return; // skip exact duplicate frames
        seen.add(key);
        frames.push({ x, y, w, h, hue: fhue });
        fhue = (fhue + 20 + Math.random() * 25) % 360;
    });
    lastHue = fhue;

    // load animations — generate hues
    const animEls = [...doc.querySelectorAll('animations > animation')];
    animations.length = 0;
    selectedAnimIdx = -1;
    let ahue = Math.random() * 360;
    animEls.forEach(el => {
        const indices = (el.getAttribute('frames') || '').trim().split(/\s+/).filter(Boolean).map(Number);
        animations.push({
            name: el.getAttribute('name') || '',
            hue: ahue,
            frameIndices: indices.filter(n => !isNaN(n)),
        });
        ahue = (ahue + 20 + Math.random() * 25) % 360;
    });
    if (animations.length) selectedAnimIdx = 0;

    rebuildAnimList();
    selectAnimation(selectedAnimIdx);
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
