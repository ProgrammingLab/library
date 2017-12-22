function getCSV(){
    let req = new XMLHttpRequest();
    req.open("get", "books.csv", true);
    req.send(null);
    req.onload = function(){
      createList(req.responseText);
    }
}

function createList(str){
  let result = [];
  let books = [];
  let tmp = str.split("\n");
  
  console.log(tmp);

  for(let i=1; i<tmp.length - 1; i++){
    result[i] = tmp[i].split(',');
    books.push({
      id: result[i][0],
      title: result[i][1],
      isbn: result[i][2],
      author: result[i][4],
      publisher: result[i][5],
      price: result[i][8],
    });
  }
  console.log(books);

  let row = document.getElementById('row-template');
  let booksList = document.querySelector('ul');

  for(let i=0; i<books.length; i++) {
    let clone = row.cloneNode(true);
    row.style = '';
    clone.querySelector('.title').innerText = books[i].title;
    clone.querySelector('.author').innerText = books[i].author;
    clone.querySelector('.publisher').innerText = books[i].publisher;
    clone.querySelector('.isbn').innerText = books[i].isbn;
    clone.querySelector('.price').innerText = books[i].price;
    clone.classList.remove('hidden');
    booksList.appendChild(clone);
  }

  row.parentNode.removeChild(row);
}

let searchInput = document.getElementById('search-input');

searchInput.addEventListener('input', (e) => {
  let booksList = document.querySelector('ul');
  let row = booksList.getElementsByTagName('li');
  for(let i=0; i<row.length; i++) {
    if(row[i].innerText.toUpperCase().includes(searchInput.value.toUpperCase())) {
      row[i].classList.remove('hidden');
    }
    else row[i].classList.add('hidden');
  }
});

getCSV();
