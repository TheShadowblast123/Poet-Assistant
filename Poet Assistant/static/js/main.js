// mousePosition = {
//     x: 0,
//     y: 0
// }
function togglePantoumEnd () {
    const toggle = document.getElementById("pantoumEnd").checked;
    const relevantLines = Array.from(document.querySelectorAll(".stanza_line"));
    const input_b = relevantLines[relevantLines.length - 3];
    const input_d = relevantLines[relevantLines.length - 1];
    if (toggle) {
        input_b.classList[0] = "line_3";
        input_b.placeholder = "3";
        input_d.classList[0] = "line_1";
        input_d.placeholder = "1";
        temp = input_b.value;
        input_b.value = input_d.value;
        input_d.value = temp;
    } else {
        input_b.classList[0] = "line_1";
        input_b.placeholder = "1";
        input_d.classList[0] = "line_3";
        input_d.placeholder = "3";
        temp = input_b.value;
        input_b.value = input_d.value;
        input_d.value = temp;
    }

}
function addPantoumStanza () {
    const relevantLines = Array.from(document.querySelectorAll(".stanza_line"));
    
    const d = parseInt(relevantLines[relevantLines.length - 5].classList[0].replace('line_', '')) + 2;
    
    const lineNumberArray = [d-3, d-1, d-2, d];
    console.log(lineNumberArray);
    const inputArray = [];
    lineNumberArray.forEach(temp => {
        const tempElement = document.createElement("input");
        tempElement.placeholder = temp;
        tempElement.type = "text";
        tempElement.className = "line_" + temp.toString() + " line "+ "stanza_line";
        tempElement.addEventListener("input", pairSync);
        inputArray.push(tempElement);
    
    });
    const newStanza = document.createElement("div");
    newStanza.className = "stanza";
    newStanza.id = "stanza";
    inputArray.forEach(el => {
        newStanza.appendChild(el);
        const sc = document.createElement("span");
        sc.className="syllable-count"
        const temp = document.createElement("br"); //if you need to read this you forgot that if you change the structure of pantoum lines, you'll have to update this code, good luck, it should be easy
        newStanza.appendChild(sc);
        newStanza.appendChild(temp);
    });
    const stanzaHolder = document.getElementById("stanza-holder");
    const lastStanza = stanzaHolder.lastElementChild;
    stanzaHolder.insertBefore(newStanza, lastStanza);
    const newRelevantLines = Array.from(document.querySelectorAll(".stanza_line"));
    const secondToLast = newRelevantLines[newRelevantLines.length-2];
    const fourthToLast = newRelevantLines[newRelevantLines.length-4];
    const temp = parseInt(newRelevantLines[newRelevantLines.length - 5].classList[0].replace('line_', ''));
    const second = temp.toString();
    const fourth = (temp - 1).toString();
    const secondToLastCorrect = "line_" + second + " line "+ "stanza_line";
    const fourthToLastCorrect = "line_" + fourth + " line "+ "stanza_line";
    fourthToLast.className = fourthToLastCorrect;
    secondToLast.className = secondToLastCorrect;
    secondToLast.placeholder = second;
    fourthToLast.placeholder = fourth;


}
function subPantoumStanza () {
    const stanzaHolder = document.getElementById("stanza-holder");
    if (stanzaHolder.children.length < 4) return;

    
    stanzaHolder.children[stanzaHolder.children.length - 2].remove();
    const relevantLines = Array.from(document.querySelectorAll(".stanza_line"));
    const secondToLast = relevantLines[relevantLines.length-2];
    const fourthToLast = relevantLines[relevantLines.length-4];
    const temp =   parseInt(relevantLines[relevantLines.length - 5].classList[0].replace('line_', ''));
    const second = temp.toString();
    const fourth = (temp - 1).toString();
    const secondToLastCorrect = "line_" + second + " line "+ "stanza_line";
    const fourthToLastCorrect = "line_" + fourth + " line "+ "stanza_line";
    fourthToLast.className = fourthToLastCorrect;
    secondToLast.className = secondToLastCorrect;
    secondToLast.placeholder = second;
    fourthToLast.placeholder = fourth;



}
// document.onmousemove = updateMousePostion;
var variableLineCount = false;
// function updateMousePostion(e) {
//     mousePosition.x = e.clientX;
//     mousePosition.y = e.clientY;
// };
var inputs;
var lines = [];
const _lines = ['', '', '']; //We probably want this to be inside of functions and not global
var endWords = ['', '', '', '', '', ''];

document.addEventListener('htmx:afterSwap', updateState);
endWords.forEach(function (n) {
    n = ' ';
});
objOfendWords = {result : endWords}
function pairSync(event) {
    if(!livePantoumToggle) return;
    const newValue = event.target.value;
    const lineClass = event.target.classList[0]; // Assuming the class structure is consistent

    document.querySelectorAll(`.${lineClass}`).forEach(input => {
        if (input !== event.target) {
            input.value = newValue;
        }
    });
}
var lineToggle, popup, popupTextElement, data_stanzas, endwordsToggle = true, livePantoumToggle = true;
var rhymeGroups; //might need this but it doesn't need to be updated
var synonymGroups; //might need this but it doesn't need to be updated
var showsLineCount;
function endWordsToggler () {
    endwordsToggle = !endwordsToggle;
}
function livePantoumToggler () {
    livePantoumToggle = !livePantoumToggle;
}
function updateState() {
    endwordsToggle, livePantoumToggle = true, true;
    if(document.getElementById("lineToggle"))
        lineToggle = document.getElementById("lineToggle");
    // if(document.getElementById("popup-text"))
    //     popup = document.getElementById("popup-text");
    // if(document.getElementById("popup-text"))
    //     popupTextElement = document.getElementById("popup-text");
    if (Array.from(document.getElementById("stanza").querySelectorAll('input')))
        inputs = Array.from(document.querySelectorAll('#stanza input'));
    data_stanzas = []
    stanzas = Array.from(document.querySelectorAll("#stanza lines"));
    
    
    // document.addEventListener("DOMContentLoaded", () => {

    //     inputs.forEach((input) => {
    //         input.addEventListener("select", () => {
    //             showPopup(input);
    //         });
    //     });
    // });
    inputs.forEach(function(input) {
        lines.push(input.value)
    });
}

updateState(); // because the server partially has the state and I don't know how and if things will load


/*   

^___ ^--- KNOWN TODOS --- ^___^
Get Popup working again
    - Send down data from server
        - Synonyms 
        - Rhymes
        - Syllables
        - Related words?
        - Needs use to find what is usefull
    - Make it pretty
        - probably should face away from the closest corner
        - probably should be an input's height away at all times
    - Make it functional
        - you should be able to copy from it without issue
        - you shouldn't be annoyed by it
        
Get multiple language support
Properly Support Sonnets

*/
function getEndWords(el, x){
    if(!endwordsToggle) return;
    x -=1
    let temp = el.value.trim().split(" ");
    if (temp[temp.length -1] != endWords[x]) el.value += ' '+ endWords[x]; 
}
function envoi(el, a, b){
    if (el.value == "") {
        el.placeholder = a.toString()+ ", " +b.toString();
    }

}
function setEndWords() {
 
    let i = 0
    const lines = document.querySelectorAll(".line")
    lines.forEach(function(line) {
        let temp = line.value.trim().split(" ");
        endWords[i] =temp[temp.length - 1];
        i++;
    });
    objOfendWords = {'result' : endWords}
}
// function hidePopup(event) {

//     const existingPopup = document.querySelector(".popup");
    
//     if (existingPopup && (!event || (event.target !== existingPopup && !existingPopup.contains(event.target)))) {
//         console.log("actually hide")
//         existingPopup.remove();
//     }
// }
function Copy() {
    text = CreatePoem();
    TextToClipboard(text);
}
function CreatePoem() {
    //needs work
    let count = 0;
    const lines  = document.querySelectorAll('.line');
    var i = 0;
    lines.forEach(line => {
        _lines[i] = line.value;
        i++;
    });
    let stanzas = Array.from(document.querySelectorAll("#stanza input"));
    // THIS NEEEDS EDITS

    var s = '\t' + _lines[0] + ' \n'  + '\t' + _lines[1] + ' \n';
    if(lineToggle.checked){
        for (var i = 0; i < stanzas.length; i++) {
            for (var j = 0; j < stanzas[i]; j++) {
                s+= (count+ 1) + ' ' + _lines[count + 1] + ' \n'
                count++;
            }
            s += '\n';
        }
    }
    else {
        _lines.forEach( line => {
            s += line + '\n'
        });
    }
    console.log(s);
    console.log(_lines)
    return s;
}
function handleChange(inputElement) {
    //likely a go problem
	//hidePopup();
    const inputValue = inputElement.value.trim();
    const lines = document.querySelectorAll('.line');
    let i = 0;
    lines.forEach(line => {
        _lines[i] = line.value;
        i++;
    });
    const syllableCount = CountSentenceSyllables(inputValue);
    let temp = inputElement
    for (let i = 0; i <2; i ++){
        const syllableCountElement = temp.nextElementSibling;
        if (syllableCountElement == null) break;
  
        if (syllableCountElement.classList.contains('syllable-count')) {
            syllableCountElement.textContent = `${syllableCount} syllables`;
            continue;
        }
        temp = syllableCountElement;

    }

}
function DownloadFile() {d
    var output = CreatePoem();

    var link = document.createElement("a");
    var file = new Blob([output], { type: 'text/plain' });
    link.href = URL.createObjectURL(file);
    if (_lines[0].value = '')
        link.download = "".concat(poemType, " draft.txt");
    else
        link.download = "".concat(_lines[0].value, ".txt");
    link.click();
    URL.revokeObjectURL(link.href);
}
function ToEditor() {

    //a go problem
}
function CountWordSyllables(w = ''){
    //Go will send us a better count for full words but this is a good temporary solution
    w = w.toLowerCase();
    //Handle exceptions
    const twoSyllable = ['coapt','coed','coinci', 'colonel', 'cafe', 'scotia'];
    if(w.length <= 3 || w == 'preach' || w == 'preyed')
        return 1;
    if(twoSyllable.includes(w))
        return 2;
    if(w == 'serious'|| w == "worcestershire"|| w == "alias"|| w == 'acacia')
        return 3;
    if(w == 'epitome'|| w == 'hyperbole')
        return 4;

    const regex = [/[aeiouy]{1,}/g, /^(?:mc)/, /^(?:tri)[aeiouy]/,
                   /^(?:bi)[aeiouy]/, /^(?:pre)[aeiouy]/, /[^tc](?:ian)$/,
                    /[aeiou][^aeiouy]e$/,/[aeiou][^aeiouy]es$/, /i[ao]$/, /[aeiouy]ing$/];
    
    //base case
    const holdVowels = w.match(regex[0]); //count vowel groups
    let count = holdVowels ? holdVowels.length: 0;
    if(count == 0) return 0;
    if(regex[1].test(w) || regex[2].test(w) || regex[3].test(w) || regex[4].test(w)) // Handle syllablic prefixes
        count++;
    if(regex[5].test(w)) // Handle ian suffix
        count++;
    if(regex[6].test(w) || regex[7].test(w)) // handle [vowel]_e(s)
        count--;
    if(regex[8].test(w)) // handle i[ao]
        count++;
    if(regex[9].test(w)) // handle vowel+ing
        count++;
    return count;
}
function CountSentenceSyllables(s ='') {
    //probably a go problem
    const wordsOnly = /\S+/g;
    const hold = s.match(wordsOnly)
    let sum = 0;
    if(hold)
        hold.forEach(word => {
            sum += CountWordSyllables(word);
        });
    return sum;
}
function TextToClipboard(text){
    if (navigator.clipboard) {
        navigator.clipboard.writeText(text)
        .then(() => {
            alert('Your Poem has been copied');
        })
        .catch(err => {
            console.error('Unable to copy text to clipboard', err)
            fallbackCopyTextToClipboard(text);
        });
    } else {
        // Fallback to execCommand if Clipboard API is not supported
        fallbackCopyTextToClipboard(text);
    }
}


function fallbackCopyTextToClipboard(text) {
    // Create a temporary input element
    var input = document.createElement('input');

    // Set its value to the text you want to copy
    input.value = text;

    // Append it to the document
    document.body.appendChild(input);

    // Select the text in the input
    input.select();

    try {
        // Execute the "copy" command to copy the selected text to the clipboard
        document.execCommand('copy');
        console.log('Text copied to clipboard (fallback)');
    } catch (err) {
        console.error('Unable to copy text to clipboard', err);
        alert('Could not copy :C');
    } finally {
        // Remove the temporary input element
        document.body.removeChild(input);
    }
}
function getSelectedWord(input) {
	
    const inputValue = input.value;
    const inputStartIndex = input.selectionStart;
    const inputEndIndex = input.selectionEnd -1;
    let startIdx = inputStartIndex;
    let endIdx = inputEndIndex;
    
    // Find the start index of the selected word
    while (startIdx > 0 && inputValue[startIdx - 1] !== ' ') {
        startIdx--;
    }
    // Find the end index of the selected word
    while (endIdx < inputValue.length && inputValue[endIdx] !== ' ') {
        endIdx++;
    }
    // Construct the selected word by extracting the substring from start to end index
    const selectedWord = inputValue.substring(startIdx, endIdx).trim();
    return selectedWord.toLowerCase();
}
function getSynonym(word) {
	//definitely a go problem

}
// function showPopup(input) {
	
//     const selectedWord = getSelectedWord(input);
//     if (!selectedWord) {
//         hidePopup();
//         return;
//     }

//     hidePopup();
    
//     popup.style.left = `${mousePosition.x}px`;
//     popup.style.top = `${mousePosition.y}px`;
    

//     document.addEventListener("click", hidePopup);

//     // Clean up the event listener when the pop-up is removed

// }


