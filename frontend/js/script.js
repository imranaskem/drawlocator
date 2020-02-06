if (window.location.protocol === 'https:') {
  var apiURL = 'https://' + window.location.host + '/api/staff'
  var wsURL = 'wss://' + window.location.host + '/websocket'
  var locationURL = 'https://' + window.location.host + '/api/locations'
} else {
  var apiURL = 'http://' + window.location.host + '/api/staff'
  var wsURL = 'ws://' + window.location.host + '/websocket'
  var locationURL = 'http://' + window.location.host + '/api/locations'
}

new Vue({
  el: "#app",

  data: {
    title: `DrawLocator`,
    search: ``,
    people: [],
    places: [],
  },

  created: function() {
      c = new WebSocket(wsURL);
      var self = this
      c.onmessage = function (event) {
        self.people = JSON.parse(event.data);
      }
      this.getLocations()
  },

  computed: {
    filteredPeople: function() {
      var self = this
      var term = this.search.toLowerCase()
      var filtered = this.people.filter(p => p.firstname.toLowerCase().includes(term) ||
                                        p.lastname.toLowerCase().includes(term) ||
                                        p.placeofwork.toLowerCase().includes(term))
      return filtered
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
    },

    getLocations: function() {
      var xhr = new XMLHttpRequest()
      var self = this
      xhr.open('GET', locationURL)
      xhr.onload = function() {
        self.places = JSON.parse(xhr.responseText)
      }
      xhr.send(null)
    }
 },

 template: `
 <div>
  <h1>{{title}}</h1>
    <input type="text" placeholder="Search..." v-model="search" />
          <table>
            <tr>
              <th>Name</th>
              <th>Location</th>
              <th>New Location</th>
            </tr>
            <tr v-for="person in filteredPeople">
              <td id="name">{{person.firstname}} {{person.lastname}}</td>
              <td id="location">{{person.placeofwork}}</td>
              <td>
                <select required v-on:change="selectChange($event, person.id)">
                  <option hidden>Please select</option>
                  <option v-for="place in places">{{place}}</option>
                </select>
              </td>
            </tr>
          </table>
  </div>
 `
})