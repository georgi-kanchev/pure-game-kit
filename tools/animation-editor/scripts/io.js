// Save/Load stubs — to be wired up later
document.getElementById('save').addEventListener('click', () => {
    // TODO: save animation data
});

document.getElementById('load').addEventListener('click', () => {
    // TODO: load animation data
});

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
