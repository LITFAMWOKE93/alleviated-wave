
import { mat4 } from '../matrix-gl/gl-matrix/dist/esm/index.js';

/// grab the canvas



document.addEventListener('DOMContentLoaded', ()=> {
    const canvas = document.getElementById("gl-canvas");
    // Set canvas size to browser client size
    canvas.width = canvas.clientWidth;
    canvas.height = canvas.clientHeight;
    

    let triangles = [
        {  // the Z coordinate determintes the depth
            positions: [
                -0.5, -0.5, -0.5,
                0.5, -0.5, -0.5,
                0.0,  0.5, -0.5,
            ],
            colors: [
                0.0, 1.0, 1.0, 1.0,    // Cyan
                0.0, 1.0, 1.0, 1.0,
                0.0, 1.0, 1.0, 1.0,
            ],
            depth: -0.5,
        },
        { 
            positions: [
                0.0, -0.5, -0.5,
                1.0, -0.5, -0.5,
                0.5,  0.5, -0.5,
           ],
            colors: [
                1.0, 1.0, 0.0, 1.0,    // Yellow
                1.0, 1.0, 0.0, 1.0,
                1.0, 1.0, 0.0, 1.0,
            ],
            depth: -0.25,
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

    let triangleBuffers = triangles.map(triangle => initBuffers(gl, triangle));
    
    // Priming draw call
    initDraw();
    // Painters algo pseudocode
    // sort polygons by their depth values
    // for each polygon:
      // for each pixel that p covers:
       // draw p.color on pixel

    // However the broad view painters algorithm simply decides which full polygon to draw first which would be the background values


 

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
// Change depth value of triangle
function updateTriangleDepth(index, newDepth) {
    triangles[index].depth = newDepth;

   
    // this changes the z value for the vertices
    // Stride over the vertice array by threes, we are only walking on the z value
   for (let i = 2; i < triangles[index].positions.length; i += 3) {
        // set z value to the new depth value
        triangles[index].positions[i] = newDepth;
   }

   


   

}

function render(gl, programInfo, buffers) {

    // Bind the position buffer and set vertex attributes
    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.position);
    gl.vertexAttribPointer(programInfo.attribLocations.vertexPosition, 3, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(programInfo.attribLocations.vertexPosition);

    // Bind the color buffer and set color attributes
    gl.bindBuffer(gl.ARRAY_BUFFER, buffers.color);
    gl.vertexAttribPointer(programInfo.attribLocations.vertexColor, 4, gl.FLOAT, false, 0, 0);
    gl.enableVertexAttribArray(programInfo.attribLocations.vertexColor);

    // Use the shader program
    gl.useProgram(programInfo.program);

    // Draw the object
    gl.drawArrays(gl.TRIANGLES, 0, buffers.vertexCount);
}




function initDraw() {

    gl.clearColor(0.0, 0.0, 0.0, 1.0); // Clear to black, fully opaque
    gl.clearDepth(1.0);                // Clear everything
    gl.enable(gl.DEPTH_TEST);          // Enable depth testing
    gl.depthFunc(gl.LEQUAL);           // Near things obscure far things
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);


    // Depth sort before drawing
    triangles.sort((a,b) => a.depth - b.depth)
    
    // This 
    triangleBuffers = triangles.map(triangle => initBuffers(gl, triangle));

    triangleBuffers.forEach((buffers, index) => {
        render(gl, programInfo, buffers, index); 
    });

    updateInfo();
}

// Event handling

document.getElementById("triangle0Depth").addEventListener('input', (event) => {
    updateTriangleDepth(0, parseFloat(event.target.value));

    initDraw();
});

document.getElementById("triangle1Depth").addEventListener('input', (event) => {
    updateTriangleDepth(1, parseFloat(event.target.value));

    initDraw();
});




function updateInfo() {
    const infoContainer = document.getElementById('triangleInfo');
    // Clear previous content
    infoContainer.innerHTML = '';

    // Append new content based on the current state
    triangles.forEach((triangle, index) => {
        const info = document.createElement('p');
        // Use template literals for dynamic values
        info.textContent = `Triangle ${index}: Depth = ${triangle.depth}`;
        infoContainer.appendChild(info); // Corrected from info.Container to infoContainer
    });
}
});

