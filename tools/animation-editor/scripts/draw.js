function getViewSize() {
    if (image) return { w: image.width, h: image.height };
    const base = 512;
    return { w: base, h: base };
}

function buildChecker(w, h) {
    checkerCanvas = document.createElement('canvas');
    checkerCanvas.width = w;
    checkerCanvas.height = h;
    const cctx = checkerCanvas.getContext('2d');
    cctx.fillStyle = '#222222';
    cctx.fillRect(0, 0, w, h);
    cctx.fillStyle = '#2a2a2a';
    for (let y = 0; y < h; y += gridSize) {
        const rowEven = Math.floor(y / gridSize) % 2 === 0;
        for (let x = rowEven ? 0 : gridSize; x < w; x += gridSize * 2) {
            cctx.fillRect(x, y, gridSize, gridSize);
        }
    }
}

function drawView() {
    const { w, h } = getViewSize();

    ctx.save();
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.translate(camera.x, camera.y);
    ctx.scale(camera.zoom, camera.zoom);

    if (image) {
        ctx.imageSmoothingEnabled = false;
        ctx.drawImage(checkerCanvas, 0, 0);
        ctx.drawImage(image, 0, 0);

        // frames
        frames.forEach((f, i) => {
            ctx.fillStyle = `hsla(${f.hue}, 55%, 50%, 0.2)`;
            ctx.fillRect(f.x, f.y, f.w, f.h);
            ctx.strokeStyle = `hsla(${f.hue}, 55%, 50%, 0.8)`;
            ctx.lineWidth = 1.5 / camera.zoom;
            ctx.strokeRect(f.x, f.y, f.w, f.h);

            const fontSize = Math.min(f.h * 0.3, f.w * 0.25, 18 / camera.zoom);
            const pad = 4 / camera.zoom;
            const label = String(i + 1);
            ctx.font = `bold ${fontSize}px 'Segoe UI', sans-serif`;
            ctx.textAlign = 'left';
            ctx.textBaseline = 'top';
            ctx.strokeStyle = '#000';
            ctx.lineWidth = 3 / camera.zoom;
            ctx.strokeText(label, f.x + pad, f.y + pad);
            ctx.fillStyle = '#fff';
            ctx.fillText(label, f.x + pad, f.y + pad);

            ctx.textBaseline = 'alphabetic';
        });

        // selection overlay
        if (selection) {
            ctx.fillStyle = 'rgba(255, 157, 92, 0.15)';
            ctx.fillRect(selection.x, selection.y, selection.w, selection.h);
            ctx.strokeStyle = 'rgba(255, 157, 92, 0.8)';
            ctx.lineWidth = 1 / camera.zoom;
            ctx.strokeRect(selection.x, selection.y, selection.w, selection.h);

            // pending frame index while Enter is held
            if (enterDigits) {
                const fontSize = Math.min(selection.h * 0.3, selection.w * 0.3, 18 / camera.zoom);
                ctx.font = `bold ${fontSize}px 'Segoe UI', sans-serif`;
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.strokeStyle = '#000';
                ctx.lineWidth = 3 / camera.zoom;
                ctx.strokeText(enterDigits, selection.x + selection.w / 2, selection.y + selection.h / 2);
                ctx.fillStyle = '#fff';
                ctx.fillText(enterDigits, selection.x + selection.w / 2, selection.y + selection.h / 2);
                ctx.textAlign = 'left';
                ctx.textBaseline = 'alphabetic';
            }
        }
    } else {
        ctx.fillStyle = '#1a1a1a';
        ctx.fillRect(0, 0, w, h);

        ctx.strokeStyle = 'rgba(255,255,255,0.1)';
        ctx.lineWidth = 1 / camera.zoom;
        ctx.strokeRect(0, 0, w, h);

        ctx.strokeStyle = 'rgba(255,255,255,0.06)';
        ctx.beginPath();
        ctx.moveTo(w / 2, 0);
        ctx.lineTo(w / 2, h);
        ctx.moveTo(0, h / 2);
        ctx.lineTo(w, h / 2);
        ctx.stroke();
    }

    ctx.restore();
}

function resetView() {
    const { w, h } = getViewSize();
    const editorRect = document.querySelector('.editor-view').getBoundingClientRect();
    const canvasRect = canvas.getBoundingClientRect();
    camera.zoom = Math.min(editorRect.width / w, editorRect.height / h) * 0.85;
    camera.x = (editorRect.left - canvasRect.left) + (editorRect.width - w * camera.zoom) / 2;
    camera.y = (editorRect.top - canvasRect.top) + (editorRect.height - h * camera.zoom) / 2;
    drawView();
}

function updateSize() {
    const view = document.querySelector('.editor-view');
    canvas.width = view.clientWidth;
    canvas.height = view.clientHeight;
    drawView();
}
