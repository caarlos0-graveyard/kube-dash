Vue.use(Toasted)

var app = new Vue({
  el: '#app',
  data: {
    svcs: []
  },
  methods:{
    get() {
      axios({ method: "GET", "url": "/api/deployments" }).then(result => {
        this.svcs = result.data;
        document.querySelector('table#svcs').classList.remove('d-none');
      }, error => {
        console.error(error);
      });
    },
    change(state, ns, name) {
      axios({
        method: "PUT",
        url: `/api/deployments/${ns}/${name}/${state}`
      }).then(result => {
        this.$toasted.show(`${name} (${ns}): going ${state}`, {duration: 5000})
        this.get();
      }, error => {
        console.error(error);
      });
    }
  },
  beforeMount(){
    this.get();
  }
})
