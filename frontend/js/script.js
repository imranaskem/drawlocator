if (window.location.protocol === 'https:') {
  var apiURL = 'https://' + window.location.host + '/staff'
  var wsURL = 'wss://' + window.location.host + '/websocket'
} else {
  var apiURL = 'http://' + window.location.host + '/staff'
  var wsURL = 'ws://' + window.location.host + '/websocket'
}

new Vue({
  el: "#app",

  data: {
    title: `DrawLocator`,
    people: null,
    places: ["Weston Street", "Baker Street", "Holiday", "Sick", "Client Office", "Working from Home"],    
  },

  created: function() {
      c = new WebSocket(wsURL);
      var self = this
      c.onmessage = function (event) {
        self.people = JSON.parse(event.data);
      }
  },

  methods: {
    selectChange: function (event, id) {
      var place = event.target.value
      this.update(place, id)
      var self = this
      setTimeout(function () {
        event.srcElement.options.selectedIndex = 0
      }, 1000)
    },

    update: function (place, id) {
      var xhr = new XMLHttpRequest()
      var self = this
      xhr.open('PATCH', apiURL + "/" + id)
      var patch = {placeofwork: place}
      var payload = JSON.stringify(patch)
      xhr.send(payload)
    }
 },

 template: `
 <div>
  <h1>{{title}}</h1>
          <table>
            <tr>
              <th>Name</th>
              <th>Location</th>
              <th colspan="2">New Location</th>
            </tr>
            <tr v-for="person in people">
              <td>{{person.firstname}} {{person.lastname}}</td>
              <td>{{person.placeofwork}}</td>
              <td>
                <select v-on:change="selectChange($event, person.id)">
                  <option>Please select</option>
                  <option v-for="place in places">{{place}}</option>
                </select>
              </td>
            </tr>
          </table>
  </div>
 `
})