


/// grab the canvas

document.addEventListener('DOMContentLoaded', ()=> {
    const canvas = document.getElementById("gl-canvas");
    // Set canvas size to browser client size
    canvas.width = canvas.clientWidth;
    canvas.height = canvas.clientHeight;


    const triangles = [
        {  // the Z coordinate determintes the depth
            positions: [
                -0.5, -0.5, -0.5,
                0.5, -0.5, -0.5,
                0.0,  0.5, -0.5,
            ],
            colors: [
                1.0, 1.0, 0.0, 1.0,    // Yellow
                1.0, 1.0, 0.0, 1.0,
                1.0, 1.0, 0.0, 1.0,
            ],
        },
        { 
            positions: [
                -0.25, -0.25, -0.25,
                0.25, -0.25, -0.25,
                0.0,  0.25, -0.25,
            ],
            colors: [
                0.0, 1.0, 1.0, 1.0,    // Cyan
                0.0, 1.0, 1.0, 1.0,
                0.0, 1.0, 1.0, 1.0,
            ],
        }
    ];
    

// Set context, handle error
const gl = canvas.getContext("webgl2");
if (!gl) {
    alert("Failure getting webgl2 context")
    return;
}

// init a program object 
const shaderProgram = initShaders(gl, "vertex-shader", "fragment-shader");
    const programInfo = {
        program: shaderProgram,
        attribLocations: {
            vertexPosition: gl.getAttribLocation(shaderProgram, 'aPosition'),
            vertexColor: gl.getAttribLocation(shaderProgram, 'aColor'),
        },
    };


    // The actual draw call
triangles.forEach(triangle => {
    const buffers = initBuffers(gl, triangle);
    render(gl, programInfo, buffers);

});


// Init buffers function for creating binding buffer data
function initBuffers(gl, data) {

    // If we are only working with vertices for this program and not matrices we can just use a float32Array as opposed to the flatten helper function that does type checking
    // Position
    const positionBuffer = gl.createBuffer();
    gl.bindBuffer(gl.ARRAY_BUFFER, positionBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(data.positions), gl.STATIC_DRAW);


    // Colors
    const colorBuffer = gl.createBuffer();
    gl.bindBuffer(gl.ARRAY_BUFFER, colorBuffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array(data.colors), gl.STATIC_DRAW);

    return {
        position: positionBuffer,
        color: colorBuffer,
        vertexCount: data.positions.length / 3 // for triangles 
    };

    
}

function render(gl, programInfo, buffers) {
    gl.clearColor(0.0, 0.0, 0.0, 1.0); // Clear to black, fully opaque
    gl.clearDepth(1.0);                 // Clear everything
    gl.enable(gl.DEPTH_TEST);           // Enable depth testing
    gl.depthFunc(gl.LEQUAL);            // Near things obscure far things

    // Clear the canvas before we start drawing on it.
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);


    // Bind the position buffer.
    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.position);
    gl.vertexAttribPointer(programInfo.attribLocations.vertexPosition, 3, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(programInfo.attribLocations.vertexPosition);

    // Bind the color buffer.
    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.color);
    gl.vertexAttribPointer(programInfo.attribLocations.vertexColor, 4, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(programInfo.attribLocations.vertexColor);

    // Draw the object
    gl.useProgram(programInfo.program);
    gl.drawArrays(gl.TRIANGLES, 0, buffers.vertexCount);
}
});