const groupList = document.getElementById('groupList');

function nextHue() {
    if (lastHue === null) {
        lastHue = Math.random() * 360;
    } else {
        lastHue = (lastHue + 20 + Math.random() * 25) % 360;
    }
    return lastHue;
}

function createGroupItem(group, idx) {
    const item = document.createElement('div');
    item.className = 'group-item';
    item.dataset.index = idx;
    item.style.setProperty('--item-color', `hsl(${group.hue}, 55%, 50%)`);

    const handle = document.createElement('span');
    handle.className = 'drag-handle';
    handle.draggable = true;
    for (let i = 0; i < 6; i++) {
        const dot = document.createElement('span');
        dot.className = 'dot';
        handle.appendChild(dot);
    }
    handle.addEventListener('dragstart', (e) => {
        e.dataTransfer.effectAllowed = 'move';
        e.dataTransfer.setData('text/plain', idx);
        item.classList.add('dragging');
    });
    handle.addEventListener('dragend', () => {
        item.classList.remove('dragging');
        groupList.querySelectorAll('.group-item').forEach(el => el.classList.remove('drag-over'));
    });

    const nameInput = document.createElement('input');
    nameInput.className = 'group-name-input';
    nameInput.value = group.name;
    nameInput.addEventListener('input', () => {
        groups[idx].name = nameInput.value;
    });
    nameInput.addEventListener('click', (e) => {
        e.stopPropagation();
        selectGroup(idx);
    });

    const cropsInput = document.createElement('input');
    cropsInput.className = 'group-crops-input';
    cropsInput.value = group.cropIndices.join(' ');
    cropsInput.placeholder = '0 1 2…';
    cropsInput.addEventListener('input', () => {
        const parts = cropsInput.value.trim().split(/\s+/).filter(Boolean);
        groups[idx].cropIndices = parts
            .map(p => parseInt(p))
            .filter(n => !isNaN(n) && n >= 0);
        drawView();
    });
    cropsInput.addEventListener('click', (e) => {
        e.stopPropagation();
        selectGroup(idx);
    });

    const delBtn = document.createElement('button');
    delBtn.className = 'group-del-btn';
    delBtn.textContent = '×';
    delBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        groups.splice(idx, 1);
        if (selectedGroupIdx >= groups.length) selectedGroupIdx = groups.length - 1;
        rebuildGroupList();
        drawView();
    });

    item.appendChild(handle);
    item.appendChild(nameInput);
    item.appendChild(cropsInput);
    item.appendChild(delBtn);

    item.addEventListener('click', () => selectGroup(idx));
    return item;
}

function rebuildGroupList() {
    groupList.innerHTML = '';
    groups.forEach((group, i) => {
        groupList.appendChild(createGroupItem(group, i));
    });
    highlightSelection();
}

function highlightSelection() {
    groupList.querySelectorAll('.group-item').forEach((el, i) => {
        el.classList.toggle('selected', i === selectedGroupIdx);
    });
}

// Drag-and-drop reorder
groupList.addEventListener('dragover', (e) => {
    e.preventDefault();
    const target = e.target.closest('.group-item');
    if (!target || target.classList.contains('dragging')) return;
    groupList.querySelectorAll('.group-item.drag-over').forEach(el => el.classList.remove('drag-over'));
    target.classList.add('drag-over');
});

groupList.addEventListener('dragleave', (e) => {
    const target = e.target.closest('.group-item');
    if (target && !target.contains(e.relatedTarget)) {
        target.classList.remove('drag-over');
    }
});

groupList.addEventListener('drop', (e) => {
    e.preventDefault();
    const target = e.target.closest('.group-item');
    if (!target) return;
    target.classList.remove('drag-over');
    const fromIdx = parseInt(e.dataTransfer.getData('text/plain'));
    const toIdx = parseInt(target.dataset.index);
    if (isNaN(fromIdx) || isNaN(toIdx) || fromIdx === toIdx) return;

    const [moved] = groups.splice(fromIdx, 1);
    groups.splice(toIdx, 0, moved);

    if (selectedGroupIdx === fromIdx) {
        selectedGroupIdx = toIdx;
    } else if (fromIdx < toIdx && selectedGroupIdx > fromIdx && selectedGroupIdx <= toIdx) {
        selectedGroupIdx--;
    } else if (fromIdx > toIdx && selectedGroupIdx >= toIdx && selectedGroupIdx < fromIdx) {
        selectedGroupIdx++;
    }

    rebuildGroupList();
    drawView();
});

const PREVIEW_MIN_H = 80;
const PREVIEW_DEFAULT_H = 220;

function selectGroup(idx) {
    stopPreview();
    selectedGroupIdx = idx;
    highlightSelection();
    const visible = idx !== -1;
    document.getElementById('previewPanel').style.display = visible ? '' : 'none';
    document.getElementById('previewResizer').style.display = visible ? '' : 'none';
    if (visible) {
        const panel = document.getElementById('previewPanel');
        if (!panel.style.height) {
            panel.style.height = PREVIEW_DEFAULT_H + 'px';
        }
    }
    drawView();
}

// Add group button
document.getElementById('addGroupBtn').addEventListener('click', () => {
    groups.push({
        name: `Group ${groups.length + 1}`,
        hue: nextHue(),
        cropIndices: [],
    });
    rebuildGroupList();
    selectGroup(groups.length - 1);
});

// Initialize
rebuildGroupList();

// Preview playback
let previewTimer = null;
let previewCropIdx = 0;

function clearPreview() {
    const pc = document.getElementById('previewCanvas');
    pc.width = 0;
    pc.height = 0;
    pc.style.width = '0';
    pc.style.height = '0';
    document.getElementById('previewCropLabel').textContent = '';
}

function stopPreview(full) {
    if (previewTimer) clearInterval(previewTimer);
    previewTimer = null;
    document.getElementById('previewToggle').textContent = '▶';
    if (full) {
        clearPreview();
    }
}

function showPreviewCrop(f, ci, group) {
    if (!image || !f) return;
    const pc = document.getElementById('previewCanvas');
    const pad = 1;
    pc.width = f.w + pad * 2;
    pc.height = f.h + pad * 2;
    const panel = document.getElementById('previewPanel');
    const maxW = panel.clientWidth - 12;
    const maxH = panel.clientHeight - 40;
    let scale = maxW / pc.width;
    if (maxH > 0 && pc.height * scale > maxH) {
        scale = maxH / pc.height;
    }
    pc.style.width = (pc.width * scale) + 'px';
    pc.style.height = (pc.height * scale) + 'px';
    const pctx = pc.getContext('2d');
    pctx.imageSmoothingEnabled = false;

    // local checkerboard anchored to preview origin, offset by pad for outline
    pctx.fillStyle = '#222222';
    pctx.fillRect(pad, pad, f.w, f.h);
    pctx.fillStyle = '#2a2a2a';
    for (let gy = pad; gy < pad + f.h; gy += gridSize) {
        const rowEven = Math.floor((gy - pad) / gridSize) % 2 === 0;
        for (let gx = rowEven ? pad : pad + gridSize; gx < pad + f.w; gx += gridSize * 2) {
            pctx.fillRect(gx, gy, gridSize, gridSize);
        }
    }

    pctx.drawImage(image, f.x, f.y, f.w, f.h, pad, pad, f.w, f.h);
    pctx.strokeStyle = group ? `hsla(${group.hue}, 55%, 50%, 0.8)` : '#3a3a3a';
    pctx.lineWidth = 1;
    pctx.strokeRect(0.5, 0.5, pc.width - 1, pc.height - 1);

    document.getElementById('previewCropLabel').textContent = ci !== undefined ? String(ci) : '';
}

function updatePreviewCrop() {
    const group = groups[selectedGroupIdx];
    if (!group || !group.cropIndices.length) return stopPreview(false);
    if (previewCropIdx >= group.cropIndices.length) {
        if (document.getElementById('previewLoop').checked) {
            previewCropIdx = 0;
        } else {
            return stopPreview(false);
        }
    }
    const ci = group.cropIndices[previewCropIdx];
    const f = crops[ci];
    if (f) {
        selection = { x: f.x, y: f.y, w: f.w, h: f.h };
        showPreviewCrop(f, ci, group);
    }
    drawView();
    previewCropIdx++;
}

// Stop preview when selecting a different group
const origSelectGroup = selectGroup;
selectGroup = function(idx) {
    stopPreview(true);
    origSelectGroup(idx);
    if (idx !== -1) {
        startPlayback();
    }
};

function startPlayback() {
    const group = groups[selectedGroupIdx];
    if (!group || !group.cropIndices.length) return;
    stopPreview(true);
    const firstCi = group.cropIndices[0];
    const firstC = crops[firstCi];
    if (firstC) {
        showPreviewCrop(firstC, firstCi, group);
    }
    previewCropIdx = 1;
    if (firstC) selection = { x: firstC.x, y: firstC.y, w: firstC.w, h: firstC.h };
    drawView();
    let speed = parseFloat(document.getElementById('previewSpeed').value);
    if (isNaN(speed)) speed = 8;
    if (speed > 0) {
        previewTimer = setInterval(updatePreviewCrop, 1000 / speed);
    }
    document.getElementById('previewToggle').textContent = '⏸';
}

document.getElementById('previewToggle').addEventListener('click', () => {
    if (previewTimer) {
        stopPreview(false);
    } else {
        startPlayback();
    }
});

// Update playback speed in real time
document.getElementById('previewSpeed').addEventListener('input', () => {
    document.getElementById('previewSpeedVal').textContent = document.getElementById('previewSpeed').value + ' FPS';
    if (!previewTimer) return;
    if (previewTimer) clearInterval(previewTimer);
    let speed = parseFloat(document.getElementById('previewSpeed').value);
    if (isNaN(speed)) speed = 8;
    if (speed > 0) {
        previewTimer = setInterval(updatePreviewCrop, 1000 / speed);
    }
});

// Preview panel vertical resizer
document.getElementById('previewResizer').addEventListener('mousedown', (e) => {
    e.preventDefault();
    const panel = document.getElementById('previewPanel');
    const parent = panel.parentElement;
    const startY = e.clientY;
    const startH = panel.offsetHeight;

    function onMove(ev) {
        const parentRect = parent.getBoundingClientRect();
        const maxH = parentRect.height - PREVIEW_MIN_H;
        let newH = startH + (startY - ev.clientY);
        newH = Math.max(PREVIEW_MIN_H, Math.min(maxH, newH));
        panel.style.height = newH + 'px';
    }

    function onUp() {
        document.removeEventListener('mousemove', onMove);
        document.removeEventListener('mouseup', onUp);
        document.body.style.cursor = 'default';
    }

    document.addEventListener('mousemove', onMove);
    document.addEventListener('mouseup', onUp);
    document.body.style.cursor = 'row-resize';
});
