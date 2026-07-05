let isPanning = false;
let isSelecting = false;
let selStart = null;
let lastMousePos = { x: 0, y: 0 };

function updateSelection(world) {
    const cx = snapToGrid(world.x);
    const cy = snapToGrid(world.y);
    const x1 = Math.min(selStart.x, cx);
    const y1 = Math.min(selStart.y, cy);
    const x2 = Math.max(selStart.x, cx) + gridSize;
    const y2 = Math.max(selStart.y, cy) + gridSize;
    selection = { x: x1, y: y1, w: x2 - x1, h: y2 - y1 };
}

canvas.addEventListener('mousedown', (e) => {
    if (e.button === 1) {
        isPanning = true;
    } else if (e.button === 2) {
        const world = screenToWorld(e.clientX, e.clientY);
        const idx = frames.findLastIndex(f =>
            world.x >= f.x && world.x <= f.x + f.w &&
            world.y >= f.y && world.y <= f.y + f.h);
        if (idx !== -1) {
            frames.splice(idx, 1);
            // remove reference from all animations
            animations.forEach(a => {
                a.frameIndices = a.frameIndices
                    .map(i => i > idx ? i - 1 : i)
                    .filter(i => i !== idx);
            });
            rebuildAnimList();
            drawView();
        }
    } else if (e.button === 0 && image) {
        const world = screenToWorld(e.clientX, e.clientY);
        if (selection && 
            world.x >= selection.x && world.x <= selection.x + selection.w &&
            world.y >= selection.y && world.y <= selection.y + selection.h) {
            selection = null;
            if (selectedAnimIdx !== -1) {
                stopPreview(true);
                selectedAnimIdx = -1;
                highlightSelection();
                document.getElementById('previewPanel').style.display = 'none';
            }
            drawView();
            return;
        }
        // deselect animation when clicking canvas
        if (selectedAnimIdx !== -1) {
            stopPreview(true);
            selectedAnimIdx = -1;
            highlightSelection();
            document.getElementById('previewPanel').style.display = 'none';
        }
        isSelecting = true;
        selStart = { x: snapToGrid(world.x), y: snapToGrid(world.y) };
        updateSelection(world);
        drawView();
    }
    lastMousePos = { x: e.clientX, y: e.clientY };
});

window.addEventListener('mousemove', (e) => {
    if (isPanning) {
        camera.x += e.clientX - lastMousePos.x;
        camera.y += e.clientY - lastMousePos.y;
        lastMousePos = { x: e.clientX, y: e.clientY };
        drawView();
    } else if (isSelecting) {
        const world = screenToWorld(e.clientX, e.clientY);
        updateSelection(world);
        drawView();
    }
});

canvas.addEventListener('wheel', (e) => {
    handleZoom(e);
    drawView();
}, { passive: false });

window.addEventListener('mouseup', () => {
    isPanning = false;
    isSelecting = false;
});

let isEnterHeld = false;

window.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && selection && image) {
        if (!isEnterHeld) {
            isEnterHeld = true;
            enterDigits = '';
            drawView();
        }
        e.preventDefault();
    } else if (isEnterHeld && e.key >= '0' && e.key <= '9') {
        enterDigits += e.key;
        drawView();
    }
});

window.addEventListener('keyup', (e) => {
    if (e.key === 'Enter' && isEnterHeld) {
        isEnterHeld = false;
        const idx = enterDigits !== '' ? parseInt(enterDigits) : frames.length;
        const insertAt = Math.min(idx, frames.length);
        if (lastHue === null) {
            lastHue = Math.random() * 360;
        } else {
            lastHue = (lastHue + 20 + Math.random() * 25) % 360;
        }
        frames.splice(Math.max(0, insertAt), 0, {
            x: selection.x,
            y: selection.y,
            w: selection.w,
            h: selection.h,
            hue: lastHue,
        });
        // shift indices in all animations for frames after insertion point
        animations.forEach(a => {
            a.frameIndices = a.frameIndices.map(i => i >= insertAt ? i + 1 : i);
        });
        selection = null;
        enterDigits = '';
        rebuildAnimList();
        drawView();
    }
});

window.addEventListener('resize', updateSize);
updateSize();
resetView();
