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

function selectAnimation(idx) {
    selectedAnimIdx = idx;
    highlightSelection();
    drawView();
}

// Add animation button
document.getElementById('addAnimBtn').addEventListener('click', () => {
    animations.push({
        name: `Animation ${animations.length + 1}`,
        hue: nextHue(),
        frameIndices: [],
    });
    selectedAnimIdx = animations.length - 1;
    rebuildAnimList();
    drawView();
});

// Initialize
rebuildAnimList();
