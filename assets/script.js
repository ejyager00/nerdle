var url = window.location.href;
var gameKey = -1;
var guessNum = 0;
var victory = false;
var currentMaxGuesses = 6;
window.onload = function() {
    document.getElementById("formSubmit").addEventListener("click", gameFormSubmitted);
    document.getElementById("guessSubmit").addEventListener("click", guessFormSubmitted);
    document.getElementById("userguess").addEventListener("keyup", function() {
        if (event.keyCode == 13) {
            guessFormSubmitted();
        }
    });
    gameFormSubmitted();
};

function gameFormSubmitted() {
    currentMaxGuesses = document.getElementById("maxguesses").value;
    let puzzlelength = document.getElementById("puzzlelength").value;
    let leadingzeros = document.getElementById("leadingzeros").checked;
    let negativezero = document.getElementById("negativezero").checked;
    let zeroremovalrate = document.getElementById("zeroremovalrate").value;
    newGame(
        puzzlelength,
        currentMaxGuesses,
        leadingzeros,
        negativezero,
        zeroremovalrate
    ).then(data => {
        document.getElementById("warning").innerHTML = "";
        var newTable = "";
        for (let i = 0; i < currentMaxGuesses; i++) {
            newTable += "<tr>";
            for (let j = 0; j < puzzlelength; j++) {
                newTable+='<td id="cell'+ i + j + '"></td>';
            }
            newTable += "</tr>";
        }
        document.getElementById("gameGrid").innerHTML = newTable;
        gameKey = data['key'];
        guessNum = 0;
        victory = false;
    });
}

function guessFormSubmitted() {
    if (guessNum < currentMaxGuesses && !victory) {
        let userGuess = document.getElementById("userguess").value;
        makeGuess(gameKey, userGuess).then(data => {
            if (data['validguess']) {
                document.getElementById("warning").innerHTML = "";
                document.getElementById("userguess").value = "";
                for (let i=0; i < data['comparison'].length; i++) {
                    document.getElementById("cell"+guessNum+i).innerHTML = userGuess.charAt(i);
                    if (data['comparison'][i] == 1) {
                        document.getElementById("cell"+guessNum+i).style.backgroundColor = "green";
                    } else if (data['comparison'][i] == -1) {
                        document.getElementById("cell"+guessNum+i).style.backgroundColor = "yellow";
                    } else {
                        document.getElementById("cell"+guessNum+i).style.backgroundColor = "red";
                    }
                }
                guessNum++;
                if (data['won']) {
                    if (guessNum == 1) {
                        document.getElementById("warning").innerHTML = "You won in one guess!";
                    } else {
                        document.getElementById("warning").innerHTML = "You won in " + guessNum + " guesses!";
                    }
                    victory = true;
                } else if (data['lost'] || guessNum >= currentMaxGuesses) {
                    document.getElementById("warning").innerHTML = "You lose! The solution was " + data['solution'] + ".";
                }
            } else {
                document.getElementById("warning").innerHTML = "Your guess is invalid!";
            }
        })
    }
}

async function newGame(puzzlelength, maxguesses, leadingzeros, negativezero, zeroremovalrate) {
    var data = '{"Length":'+puzzlelength+
    ',"LeadingZeros":'+leadingzeros+
    ',"NegativeZeros":'+negativezero+
    ',"ZeroRemovalRate":'+zeroremovalrate+
    ',"MaxGuesses":'+maxguesses+'}';
    const response = await fetch(url+"new", {
        method : "POST",
        mode : "cors",
        headers: {
            'Content-Type': 'application/json'
        },
        body: data
    });
    return response.json();
}
async function makeGuess(key, guess) {
    var data = JSON.stringify({
        Key:key,
        Guess:guess
    });
    const response = await fetch(url+"guess", {
        method : "POST",
        mode : "cors",
        headers: {
            'Content-Type': 'application/json'
        },
        body: data
    });
    return response.json();
}