// prevent default for the form

var xhttp = new XMLHttpRequest();

// from prevent default
const partyForm = document.getElementById('partyForm');
partyForm.addEventListener('submit', (e) => {
    e.preventDefault();
});

const editBtn = document.getElementById('editBtn');
const cancelBtn = document.getElementById('cancelBtn');
const saveBtn = document.getElementById('saveBtn');

const formExtractValues = (form, vals) =>{
    if (vals == undefined){
        var vals = {}
    }
    const elems = form.children
    for(var i=0; i < elems.length; i++){
        let elem = elems[i]
        if(elem.type == 'text'){
            vals[elem.id] = elem.value
        }else if (elem.children.length > 0){
            formExtractValues(elem, vals)
        }
    }
    return vals
}

// iterate over children until a input is found and toggle disable propoerty
const toggleDisabled = (elems) => {
    for(var i=0; i < elems.length; i++){
        let elem = elems[i]        
        if(elem.type == 'text'){
            if(elem.disabled == true){
                elem.disabled = false;
            }else{
                elem.disabled = true;
            }
        }else if (elem.children.length > 0){
            toggleDisabled(elem.children)
        }
    }
}

const switchEdit = (ev) => {
    toggleDisabled(partyForm.children)
    // set cancel button
};

const postData = () => {
    // get data from the form
    // send it to the server
    console.log('saving data')
    //form.submit()
    
    let party = formExtractValues(partyForm)
    console.log(party);
    const url = `/api/party/${party.taxpayerId}`
    xhttp.open("PUT", url, true);

    party.birthDate = `${party.year}-${party.month}-${party.day}`
    xhttp.send(JSON.stringify(party));
    // TODO: handle async
};

editBtn.addEventListener('click', switchEdit)
saveBtn.addEventListener('click', postData)
