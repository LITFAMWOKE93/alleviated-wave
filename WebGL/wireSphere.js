"use strict";

var wireSphere = function() {
    var canvas;
    var gl;

    var numTimesDivide = 6;

    var index = 0;

    var positionsArray = [];


    var near = -10;
    var far = 10;
    var radius = 6.0;
    var theta = 0.0;
    var phi = 0.0;
    var dr = 5.0* Math.PI/180.0;

    var left = -2.0;
    var right = 2.0;
    var top = 2.0;
    var bottom = -2.0;

    // model view and projection matrix
    var MVM, PM

    // atrritube locations in the shader program
    var MVMloc, PMloc;
    // storage for view
    var eye;
    // LookAt world origin
    const at = vec3(0.0,0.0,0.0);
    // Y axis is used to orient up direction
    const up = vec3(0.0,1.0,0.0);

    // Get the three vertices to form a triangle
    function triangle(a,b,c) {
        positionsArray.push(a);
        positionsArray.push(b);
        positionsArray.push(c);
        index += 3
    }
    // Recursive call for division
    function divideTriangle(a,b,c,count) {

        if (count > 0) {
            var ab = normalize(mix(a,b,0.5), true);
            var ac = normalize(mix(a,c,0.5), true);
            var bc = normalize(mix(b,c,0.5), true);

            divideTriangle(a, ab, ac, count - 1);
            divideTriangle(ab, b, bc, count - 1);
            divideTriangle(bc, c, ac, count - 1);
            divideTriangle(ab, bc, ac, count - 1);
        } else {
            // Base case for recursive method
            triangle(a,b, c);
        }
    }

    // Create tetrahedron from triangle subdivision
    function tetrahedron(a, b, c, d, n) {
        divideTriangle(a, b, c, n);
        divideTriangle(d, c, b, n);
        divideTriangle(a, d, b, n);
        divideTriangle(a, c, d, n);
    }

    window.onload = function init() {
        // Accessing the DOM to select and HTML element by ID
        canvas = document.getElementById("gl-canvas")

        gl = canvas.getContext('webgl2');
        if (!gl) alert("WebGL 2.0 isn't available");

        // Create viewport equal to the canvas defined in HTML
        gl.viewport(0, 0, canvas.width, canvas.width);
        // Clear color is set to white
        gl.clearColor(1.0,1.0,1.0,1.0);

        // Load shaders
        
        var program = initShaders(gl, "vertex-shader", "fragment-shader");
        gl.useProgram(program)

        // These four vectors define a base tetrahedron 
        var va = vec4(0.0, 0.0, -1.0, 1);
        var vb = vec4(0.0, 0.942809, 0.333333, 1);
        var vc = vec4(-0.816497, -0.471405, 0.333333, 1);
        var vd = vec4(0.816497, -0.471405, 0.333333, 1);

        tetrahedron(va,vb,vc,vd, numTimesDivide);
        // Init attribute buffers
        var vBuffer = gl.createBuffer();
        gl.bindBuffer(gl.ARRAY_BUFFER, vBuffer);
        gl.bufferData( gl.ARRAY_BUFFER, flatten(positionsArray), gl.STATIC_DRAW);
        

        var positionLoc = gl.getAttribLocation(program, "aPosition");
        gl.vertexAttribPointer( positionLoc, 4, gl.FLOAT, false, 0, 0);
        gl.enableVertexAttribArray(positionLoc);

        // Get addresses of shader variables
        MVMloc = gl.getUniformLocation(program, "uModelViewMatrix");
        PMloc = gl.getUniformLocation(program, "uProjectionMatrix");

        // Event listeners for increase/decrease angles
        document.getElementById("Button0").onclick = function(){theta += dr;};
        document.getElementById("Button1").onclick = function(){theta -= dr;};
        document.getElementById("Button2").onclick = function(){phi += dr;};
        document.getElementById("Button3").onclick = function(){phi -= dr;};
        // Event listeners for division call
        document.getElementById("Button4").onclick = function(){
            numTimesToSubdivide++;
            index = 0;
            positionsArray = [];
            init();
        };
        // Event listener decreased divison call
        document.getElementById("Button5").onclick = function() {
            if(numTimesDivide) numTimesDivide--;
            index = 0;
            positionsArray = [];
            init();
        }
        // Call render loop
        render();
    }

function render() {
    gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
    
    // viewing vector
    eye = vec3(radius*Math.sin(theta)*Math.cos(phi), radius*Math.sin(theta)*Math.sin(phi), radius*Math.cos(theta));

    // Create matrix for model-view
    MVM = lookAt(eye, at, up);
    // Create projection matrix
    PM = ortho(left, right, bottom, top, near, far);

    // give the GPU the values calculated on the CPU
    gl.uniformMatrix4fv(MVMloc, false, flatten(MVM));
    gl.uniformMatrix4fv(PMloc, false, flatten(PM));

    for (var i= 0; i<index; i+=3)
    gl.drawArrays(gl.LINE_LOOP, i, 3);

    requestAnimationFrame(render);

}



}

// Entire program wrapped in a single function call
wireSphere();

