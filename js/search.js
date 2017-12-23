function getCSV() {
    let req = new XMLHttpRequest();
    req.open("get", "books.csv", true);
    req.send(null);
    req.onload = function() {
        createList(req.responseText);
    }
}

function createList(str) {
    let result = [];
    let books = [];
    let tmp = str.split("\n");

    for (let i = 1; i < tmp.length - 1; i++) {
        result[i] = tmp[i].split(',');
        books.push({
            id: result[i][0],
            title: result[i][1],
            isbn: result[i][2],
            author: result[i][4],
            publisher: result[i][5],
        });
    }

    let row = document.getElementById('row-template');
    let booksList = document.querySelector('ul');

    for (let i = 0; i < books.length; i++) {
        let clone = row.cloneNode(true);
        row.style = '';
        clone.querySelector('.title').innerText = books[i].title;
        clone.querySelector('.author').innerText = books[i].author;
        clone.querySelector('.publisher').innerText = books[i].publisher;
        clone.querySelector('.isbn').innerText = books[i].isbn;
        clone.classList.remove('hidden');
        booksList.appendChild(clone);
    }

    row.parentNode.removeChild(row);
}

let searchInput = document.getElementById('search-input');
let booksList = document.getElementById('books-list');
let row = booksList.getElementsByTagName('li');
let displayFlag = new Array(row.length);
let image = document.getElementById('not-found');

searchInput.addEventListener('input', ()=>{
    let found = false;
    let string = searchInput.value.toUpperCase();
    for (let i = 0; i < row.length; i++) {
        if (row[i].innerText.toUpperCase().includes(string)) {
            displayFlag[i] = true;
            found = true;
        } else
            displayFlag[i] = false;
    }
    for (let i = 0; i < row.length; i++) {
        if (displayFlag[i])
            row[i].classList.remove('hidden');
        else
            row[i].classList.add('hidden');
    }
    if (found)
        image.classList.add('hidden');
    else
        image.classList.remove('hidden');
}
);

getCSV();
