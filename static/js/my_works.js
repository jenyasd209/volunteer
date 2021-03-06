function collapsible() {
    var elems = document.querySelectorAll('.collapsible');
    var instances = M.Collapsible.init(elems, {});
}

const select = document.getElementById('status');
const orders = document.getElementById('works');

window.onload = function(){
    getOrders();
};

function getOrders() {
    fetch('/my_profile/sort_works_by_'+select.value)
        .then(res => res.json())
        .then((json) => {
            orders.innerText = ``;
            if (select.value === select.options[0].value){
                renderPerformedOrders(json)
            }else if (select.value === select.options[1].value){
                renderDoneOrders(json);
            }
        })
        .catch(err => { throw err });
}

function renderPerformedOrders(performedOrders) {
    if (performedOrders === null){
        renderOrders(new NoOrders("in work").create());
    }else {
        for (let i in performedOrders) {
            let order = new CardPerformedOrder(performedOrders[i]);
            renderOrders(order.create());
        }
    }
}

function renderDoneOrders(doneOrders) {
    if (doneOrders === null){
        renderOrders(new NoOrders(select.value).create());
    }else {
        for (let i in doneOrders) {
            let order = new CardDoneOrder(doneOrders[i]);
            renderOrders(order.create());
            collapsible();
        }
    }
}

function renderOrders(child){
    orders.appendChild(child);
}

class CardDoneOrder{
    constructor(doneOrder) {
        this.doneOrder = doneOrder;
    }

    create(){
        let card = document.createElement("div");
        card.className = 'card order';

        card.appendChild(this.createContent());
        card.appendChild(this.createAction());

        return card;
    }

    createContent(){
        let card_content = new CardContent().create();
        let freelancer_comment = this.doneOrder.freelancer_comment.text;
        let btnLeaveComment = ``;
        if (freelancer_comment === ""){
            freelancer_comment = "Volunteer did not comment";
            btnLeaveComment = `<a class="btn-small" href="/orders/id${this.doneOrder.order.id}/freelancer_comment">Leave comment</a>`;
        }

        card_content.innerHTML = `<span class="card-title">
                                     <a href="/orders/id${this.doneOrder.order.id}"> ${this.doneOrder.order.title} </a>
                                  </span>
                                  <ul class="collapsible">
                                    <li>
                                      <div class="collapsible-header card-inform">
                                          <span>Customer comment</span>
                                          
                                          <div class="rating" id="${this.doneOrder.customer_comment.id}">
                                                ${renderStarsComment(this.doneOrder.customer_comment.id, this.doneOrder.customer_comment.rait)}
                                          </div>
                                      </div>
                                      <div class="collapsible-body">
                                          <p> <a href="/customers/id${this.doneOrder.order.customer.user.ID}"> ${this.doneOrder.order.customer.user.FirstName} ${this.doneOrder.order.customer.user.LastName} </a></p>
                                          <span>${this.doneOrder.customer_comment.text}</span>
                                      </div>
                                    </li>
                                    <li>
                                      <div class="collapsible-header card-inform">
                                        <span>My comment</span>
                                          
                                          <div class="rating" id="${this.doneOrder.freelancer_comment.id}">
                                              ${renderStarsComment(this.doneOrder.freelancer_comment.id, this.doneOrder.freelancer_comment.rait)}
                                          </div>
                                      </div>
                                      <div class="collapsible-body">
                                        <p>${freelancer_comment}</p>
                                        ${btnLeaveComment}
                                      </div>
                                    </li>
                                  </ul>`;
        return card_content
    }

    createAction(){
        let card_action = new CardAction().create();
        card_action.classList.add('card-inform');
        card_action.innerHTML = `<span class="center">Complete: ${formatDate(new Date(this.doneOrder.date_complete))}</span>`;

        return card_action
    }
}

class CardPerformedOrder{
    constructor(availableOrder) {
        this.performedOrder = availableOrder;
    }

    create(){
        let card = document.createElement("div");
        card.className = 'card order';

        card.appendChild(this.createContent());
        card.appendChild(this.createAction());
        return card;
    }

    createContent(){
        let card_content = new CardContent().create();
        card_content.innerHTML = `<span class="card-title">
                                     <a href="/orders/id${this.performedOrder.order.id}"> ${this.performedOrder.order.title} </a>
                                  </span>
                                  <span>${this.performedOrder.order.content}</span>`;

        return card_content
    }

    createAction(){
        let card_action = new CardAction().create();
        card_action.innerHTML = `<p>Status: <span class="green-text">${this.performedOrder.order.status.Name}</span></p>
                                <span>Customer: <a href="/customers/id${this.performedOrder.order.customer.user.ID}">${this.performedOrder.order.customer.user.FirstName} ${this.performedOrder.order.customer.user.LastName}</a></span>`;

        return card_action
    }
}

class CardContent{
    constructor(){}
    create(){
        let card_content = document.createElement("div");
        card_content.className = 'card-content';

        return card_content;
    }
}

class CardAction{
    constructor(){}
    create(){
        let card_action = document.createElement("div");
        card_action.className = 'card-action';

        return card_action;
    }
}

class NoOrders {
    constructor(statusName){
        this.text= 'No ' + statusName;
    }
    create(){
        let no_orders = document.createElement('p');
        no_orders.className = 'center-align flow-text';
        no_orders.innerText = this.text;

        return no_orders;
    }
}

function formatDate(date) {
    let days = ['Sunday','Monday','Tuesday','Wednesday','Thursday','Friday','Saturday'];
    let months = ['January','February','March','April','May','June','July','August','September','October','November','December'];

    let formatDate = date.getHours() + ':' + date.getMinutes() + ' | ';
    if (date.getDate() === new Date().getDate()){
        formatDate += `Today`;
    }else if (date.getDate() === (new Date().setDate(new Date().getDate() - 1))){
        formatDate += `Yesterday`
    }else{
        formatDate += date.getDate() + '.' + date.getMonth() + '.' + date.getFullYear();
    }
    return formatDate
}