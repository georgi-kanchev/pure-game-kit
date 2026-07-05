const canvas = document.getElementById('view');
const ctx = canvas.getContext('2d');

let image = null;
let gridSize = 16;
let selection = null; // { x, y, w, h } in world coords, grid-aligned
let checkerCanvas = null; // pre-rendered checkerboard at image dimensions
let frames = [];
let lastHue = null;
let enterDigits = ''; // digits typed while Enter is held

document.addEventListener('contextmenu', e => e.preventDefault());
