const animationList = document.getElementById('animationList');

function nextHue() {
    if (lastHue === null) {
        lastHue = Math.random() * 360;
    } else {
        lastHue = (lastHue + 20 + Math.random() * 25) % 360;
    }
    return lastHue;
}

function createAnimItem(anim, idx) {
    const item = document.createElement('div');
    item.className = 'anim-item';
    item.style.setProperty('--item-color', `hsl(${anim.hue}, 55%, 50%)`);

    const nameInput = document.createElement('input');
    nameInput.className = 'anim-name-input';
    nameInput.value = anim.name;
    nameInput.addEventListener('input', () => {
        animations[idx].name = nameInput.value;
    });
    nameInput.addEventListener('click', (e) => {
        e.stopPropagation();
        selectAnimation(idx);
    });

    const framesInput = document.createElement('input');
    framesInput.className = 'anim-frames-input';
    framesInput.value = anim.frameIndices.join(' ');
    framesInput.placeholder = '0 1 2…';
    framesInput.addEventListener('input', () => {
        const parts = framesInput.value.trim().split(/\s+/).filter(Boolean);
        animations[idx].frameIndices = parts
            .map(p => parseInt(p))
            .filter(n => !isNaN(n) && n >= 0);
        drawView();
    });
    framesInput.addEventListener('click', (e) => {
        e.stopPropagation();
        selectAnimation(idx);
    });

    const delBtn = document.createElement('button');
    delBtn.className = 'anim-del-btn';
    delBtn.textContent = '×';
    delBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        animations.splice(idx, 1);
        if (selectedAnimIdx >= animations.length) selectedAnimIdx = animations.length - 1;
        rebuildAnimList();
        drawView();
    });

    item.appendChild(nameInput);
    item.appendChild(framesInput);
    item.appendChild(delBtn);

    item.addEventListener('click', () => selectAnimation(idx));
    return item;
}

function rebuildAnimList() {
    animationList.innerHTML = '';
    animations.forEach((anim, i) => {
        animationList.appendChild(createAnimItem(anim, i));
    });
    highlightSelection();
}

function highlightSelection() {
    animationList.querySelectorAll('.anim-item').forEach((el, i) => {
        el.classList.toggle('selected', i === selectedAnimIdx);
    });
}

const PREVIEW_MIN_H = 80;
const PREVIEW_DEFAULT_H = 220;

function selectAnimation(idx) {
    stopPreview();
    selectedAnimIdx = idx;
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

// Add animation button
document.getElementById('addAnimBtn').addEventListener('click', () => {
    animations.push({
        name: `Animation ${animations.length + 1}`,
        hue: nextHue(),
        frameIndices: [],
    });
    rebuildAnimList();
    selectAnimation(animations.length - 1);
});

// Initialize
rebuildAnimList();

// Preview playback
let previewTimer = null;
let previewFrameIdx = 0;

function clearPreview() {
    const pc = document.getElementById('previewCanvas');
    pc.width = 0;
    pc.height = 0;
    pc.style.width = '0';
    pc.style.height = '0';
}

function stopPreview(full) {
    if (previewTimer) clearInterval(previewTimer);
    previewTimer = null;
    document.getElementById('previewFrame').disabled = false;
    if (full) {
        clearPreview();
    }
}

function showPreviewFrame(f) {
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
    pctx.drawImage(image, f.x, f.y, f.w, f.h, pad, pad, f.w, f.h);
    pctx.strokeStyle = '#3a3a3a';
    pctx.lineWidth = 1;
    pctx.strokeRect(0.5, 0.5, pc.width - 1, pc.height - 1);
}

function updatePreviewFrame() {
    const anim = animations[selectedAnimIdx];
    if (!anim || !anim.frameIndices.length) return stopPreview(false);
    if (previewFrameIdx >= anim.frameIndices.length) {
        if (document.getElementById('previewLoop').checked) {
            previewFrameIdx = 0;
        } else {
            return stopPreview(false);
        }
    }
    const fi = anim.frameIndices[previewFrameIdx];
    const f = frames[fi];
    if (f) {
        selection = { x: f.x, y: f.y, w: f.w, h: f.h };
        showPreviewFrame(f);
    }
    document.getElementById('previewFrame').value = fi;
    drawView();
    previewFrameIdx++;
}

// Stop preview when selecting a different animation
const origSelectAnimation = selectAnimation;
selectAnimation = function(idx) {
    stopPreview(true);
    origSelectAnimation(idx);
    if (idx !== -1) {
        startPlayback();
    }
};

function startPlayback() {
    const anim = animations[selectedAnimIdx];
    if (!anim || !anim.frameIndices.length) return;
    stopPreview(true);
    const firstFi = anim.frameIndices[0];
    const firstF = frames[firstFi];
    if (firstF) {
        showPreviewFrame(firstF);
        document.getElementById('previewFrame').value = firstFi;
    }
    previewFrameIdx = 1;
    if (firstF) selection = { x: firstF.x, y: firstF.y, w: firstF.w, h: firstF.h };
    document.getElementById('previewFrame').disabled = true;
    drawView();
    const speed = parseInt(document.getElementById('previewSpeed').value) || 8;
    previewTimer = setInterval(updatePreviewFrame, 1000 / speed);
}

document.getElementById('previewPlay').addEventListener('click', startPlayback);

document.getElementById('previewPause').addEventListener('click', () => {
    stopPreview(false);
});

// Update playback speed in real time
document.getElementById('previewSpeed').addEventListener('input', () => {
    if (!previewTimer) return;
    if (previewTimer) clearInterval(previewTimer);
    const speed = parseInt(document.getElementById('previewSpeed').value) || 8;
    previewTimer = setInterval(updatePreviewFrame, 1000 / speed);
});

// Frame field: type a frame index to preview it
document.getElementById('previewFrame').addEventListener('input', () => {
    const fi = parseInt(document.getElementById('previewFrame').value);
    if (isNaN(fi) || fi < 0) return;
    const f = frames[fi];
    if (f) {
        showPreviewFrame(f);
        selection = { x: f.x, y: f.y, w: f.w, h: f.h };
        drawView();
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
