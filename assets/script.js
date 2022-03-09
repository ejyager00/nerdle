const my_yellow = "#F9F786";
const my_red = "#C47482";
const my_green = "#98D4BB";
const my_white = "#F5F3E7";
const red_square = "&#128997;";
const yellow_square = "&#129000;";
const green_square = "&#129001;";
var url = window.location.href;
var gameKey = -1;
var guessNum = 0;
var victory = false;
var currentMaxGuesses = 6;
var puzzlelength = 8;
var userGuess = "";
var buttons = '0123456789+-*/=';
var input_position = 0;
var results_array = [];
var char_colors = {};
for (let i in buttons) {
    char_colors[buttons.charAt(i)] = 0;
}

window.onload = setup;

function setup() {
    // window sizing
    window.addEventListener('resize', sizeChanged);
    // Hiding elements
    document.getElementById("warning").style.display = "none";
    document.getElementById("formSubmit1").style.display = "none";
    // adding button event listeners
    document.getElementById("formSubmit1").addEventListener("click", gameFormSubmitted);
    document.getElementById("formSubmit2").addEventListener("click", gameFormSubmitted);
    document.getElementById("guessSubmit").addEventListener("click", guessFormSubmitted);
    document.getElementById("backspace").addEventListener("click", backspace);
    for (let i in buttons) {
        (function(index) {
            document.getElementById(buttons.charAt(index)).addEventListener("click", function () {
                handleCharClick(buttons.charAt(index));
            });
        })(i);
    }
    // adding key event listeners
    document.onkeyup = function(e) {
        var selectedInput = document.activeElement.id;
        if (selectedInput == "puzzlelength" || selectedInput == "maxguesses" || selectedInput=="zeroremovalrate") {
            if (e.key=="Enter") {
                gameFormSubmitted();
            }
        } else {
            if (buttons.includes(e.key)) {
                document.getElementById(e.key).click();
            } else if (e.key=="Backspace" || e.key=="Delete") {
                backspace();
            } else if (e.key=="Enter") {
                guessFormSubmitted();
            }
        }
    }
    document.addEventListener('paste', handlePaste);
    // initializing game
    gameFormSubmitted();
}

function gameFormSubmitted() {
    document.getElementsByTagName("article")[0].classList.add("no_click");
    for (let i = 0; i< document.getElementsByTagName("button").length; i++) {
        document.getElementsByTagName("button")[i].style.cursor = "wait";
        document.getElementsByTagName("button")[i].disabled = true;
    }
    for (let i = 0; i< document.getElementsByTagName("input").length; i++) {
        document.getElementsByTagName("input")[i].style.cursor = "wait";
        document.getElementsByTagName("input")[i].disabled = true;
    }
    currentMaxGuesses = document.getElementById("maxguesses").value;
    puzzlelength = document.getElementById("puzzlelength").value;
    let leadingzeros = document.getElementById("leadingzeros").checked;
    let negativezero = document.getElementById("negativezero").checked;
    let zeroremovalrate = document.getElementById("zeroremovalrate").value;
    for (let i in buttons) {
        document.getElementById(buttons.charAt(i)).style.backgroundColor = my_white;
        char_colors[buttons.charAt(i)] = 0;
    }
    newGame(
        puzzlelength,
        currentMaxGuesses,
        leadingzeros,
        negativezero,
        zeroremovalrate
    ).then(data => {
        document.getElementById("warning").innerHTML = "";
        document.getElementById("warning").style.display = "none";
        var newTable = "";
        for (let i = 0; i < currentMaxGuesses; i++) {
            newTable += "<tr>";
            for (let j = 0; j < puzzlelength; j++) {
                newTable += '<td id="cell' + i + j + '"></td>';
            }
            newTable += "</tr>";
        }
        document.getElementById("gameGrid").innerHTML = newTable;
        gameKey = data['key'];
        guessNum = 0;
        victory = false;
        userGuess = "";
        input_position = 0;
        results_array = [];
        char_colors = {};
        for (let i in buttons) {
            char_colors[buttons.charAt(i)] = 0;
        }
        document.getElementsByTagName("article")[0].classList.remove("no_click");
        for (let i = 0; i< document.getElementsByTagName("button").length; i++) {
            document.getElementsByTagName("button")[i].style.cursor = "pointer";
            document.getElementsByTagName("button")[i].disabled = false;
        }
        for (let i = 0; i< document.getElementsByTagName("input").length; i++) {
            document.getElementsByTagName("input")[i].style.cursor = "";
            document.getElementsByTagName("input")[i].disabled = false;
        }
    });
    document.activeElement.blur();
}

function guessFormSubmitted() {
    if (guessNum < currentMaxGuesses && !victory) {
        makeGuess(gameKey, userGuess).then(data => {
            console.log(data['comparison']);
            if (data['validguess']) {
                document.getElementById("warning").innerHTML = "";
                document.getElementById("warning").style.display = "none";
                for (let i = 0; i < data['comparison'].length; i++) {
                    if (data['comparison'][i] == 1) {
                        document.getElementById("cell" + guessNum + i).style.backgroundColor = my_green;
                        document.getElementById(userGuess.charAt(i)).style.backgroundColor = my_green;
                        char_colors[userGuess.charAt(i)] = 3;
                    } else if (data['comparison'][i] == -1) {
                        document.getElementById("cell" + guessNum + i).style.backgroundColor = my_yellow;
                        if (char_colors[userGuess.charAt(i)]<2) {
                            document.getElementById(userGuess.charAt(i)).style.backgroundColor = my_yellow;
                            char_colors[userGuess.charAt(i)] = 2;
                        }
                    } else {
                        document.getElementById("cell" + guessNum + i).style.backgroundColor = my_red;
                        if (char_colors[userGuess.charAt(i)]<1) {
                            document.getElementById(userGuess.charAt(i)).style.backgroundColor = my_red;
                            char_colors[userGuess.charAt(i)] = 1;
                        }
                    }
                }
                results_array.push(data['comparison']);
                guessNum++;
                userGuess = "";
                input_position = 0;
                if (data['won']) {
                    var winmessage = "";
                    if (guessNum == 1) {
                        winmessage += "You won in one guess! Copy and paste to share your results:<br>" + createSharingString(results_array);
                    } else {
                        winmessage += "You won in " + guessNum + " guesses! Copy and paste to share your results:<br>" + createSharingString(results_array);
                    }
                    document.getElementById("warning").style.display = "block";
                    document.getElementById("warning").innerHTML = winmessage;
                    victory = true;
                    document.getElementById("formSubmit1").style.display = "block";
                } else if (data['lost'] || guessNum >= currentMaxGuesses) {
                    document.getElementById("warning").innerHTML = "You lose! The solution was " + data['solution'] + ". Copy and paste to share your results:<br>" + createSharingString(results_array);
                    document.getElementById("warning").style.display = "block";
                    document.getElementById("formSubmit1").style.display = "block";
                }
            } else {
                document.getElementById("warning").innerHTML = "Your guess is invalid!";
                document.getElementById("warning").style.display = "block";
            }
        })
    }
}

function handleCharClick(c) {
    if (input_position < puzzlelength) {
        if (c=="*") {
            document.getElementById("cell" + guessNum + input_position).innerHTML = '&#215;';
        } else if (c=="/") {
            document.getElementById("cell" + guessNum + input_position).innerHTML = '&#247;';
        } else {
            document.getElementById("cell" + guessNum + input_position).innerHTML = c;
        }
        userGuess += c;
        input_position++;
    }
}

function handlePaste(e) {
    var selectedInput = document.activeElement.id;
    if (!(selectedInput == "puzzlelength" || selectedInput == "maxguesses" || selectedInput=="zeroremovalrate")) {
        var clipText = (e.clipboardData || window.clipboardData).getData('text');
        for (let i = 0; i < Math.min(puzzlelength, clipText.length); i++) {
            if (buttons.includes(clipText.charAt(i))){
                document.getElementById(clipText.charAt(i)).click();
            } else if (clipText.charAt(i)=='×') {
                document.getElementById('*').click();
            } else if (clipText.charAt(i)=='÷') {
                document.getElementById('/').click();
            }
        }
        e.preventDefault();
    }
}

function backspace() {
    if (input_position != 0) {
        input_position--;
        document.getElementById("cell" + guessNum + input_position).innerHTML = "";
        userGuess = userGuess.slice(0,-1);
    }
}

function createSharingString(results) {
    output = ""
    for (let i in results) {
        for (let j in results[i]) {
            if (results[i][j]==1) {
                output += green_square;
            } else if (results[i][j]==0) {
                output += red_square;
            } else {
                output += yellow_square;
            }
        }
        output += "<br>";
    }
    return output;
}

function sizeChanged() {
    return;
}

async function newGame(puzzlelength, maxguesses, leadingzeros, negativezero, zeroremovalrate) {
    var data = '{"Length":' + puzzlelength +
        ',"LeadingZeros":' + leadingzeros +
        ',"NegativeZeros":' + negativezero +
        ',"ZeroRemovalRate":' + zeroremovalrate +
        ',"MaxGuesses":' + maxguesses + '}';
    const response = await fetch(url + "new", {
        method: "POST",
        mode: "cors",
        headers: {
            'Content-Type': 'application/json'
        },
        body: data
    });
    return response.json();
}

async function makeGuess(key, guess) {
    var data = JSON.stringify({
        Key: key,
        Guess: guess
    });
    const response = await fetch(url + "guess", {
        method: "POST",
        mode: "cors",
        headers: {
            'Content-Type': 'application/json'
        },
        body: data
    });
    return response.json();
}