const specializations = document.getElementById('specializations');

function createSpecialization(){
    let name = document.getElementById('new-spec');
    let url = "/moderator/create_spec";
    let spec = {
        'id':0,
        'name': name.value,
    };
    fetch(url,
        {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(spec)
        })
        .then(res => res.json())
        .then((spec) => {
            let new_spec = `<div class="row white z-depth-2" id="spec-id${spec.id}">
                            <div class="col s12 input-field">
                                <input value="${spec.name}" id="spec-id${spec.id}-name">
                                <button class="btn-small right" onclick="deleteSpecialization(${spec.id})" title="Delete"><i class="material-icons">delete</i></button>
                                <button class="btn-small" onclick="editSpecialization(${spec.id})" title="Edit"><i class="material-icons">edit</i></button>
                            </div>
                        </div>`;
            name.value = "";
            specializations.innerHTML = new_spec + specializations.innerHTML;
            alert( "Specialization successful created!");
        })
        .catch(function(res){ console.log(res) });
}

function editSpecialization(id) {
    let name = document.getElementById('spec-id'+id+'-name');
    let url = "/moderator/edit_spec_id"+id;
    let spec = {
        'id':id,
        'name': name.value,
    };
    console.log(spec.name);
    fetch(url,
        {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(spec)
        })
        .catch(function(res){ console.log(res) });
    alert( "Specialization successful updated!");
}

function deleteSpecialization(id) {
    let spec = {
        'id':id,
        'name': "",
    };
    let url = "/moderator/delete_spec_id"+id;
    fetch(url,
        {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(spec)
        })
        .catch(function(res){ console.log(res) });
    document.getElementById("spec-id"+id).remove();
    alert( "Specialization successful removed!");
}