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
            drawView();
        }
    } else if (e.button === 0 && image) {
        const world = screenToWorld(e.clientX, e.clientY);
        // click on existing selection → clear it
        if (selection && 
            world.x >= selection.x && world.x <= selection.x + selection.w &&
            world.y >= selection.y && world.y <= selection.y + selection.h) {
            selection = null;
            drawView();
            return;
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
        const idx = parseInt(enterDigits) || frames.length + 1;
        const insertAt = Math.min(idx - 1, frames.length);
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
        selection = null;
        enterDigits = '';
        drawView();
    }
});

window.addEventListener('resize', updateSize);
updateSize();
resetView();
