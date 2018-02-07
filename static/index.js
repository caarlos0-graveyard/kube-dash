var app = new Vue({
  el: '#app',
  data: {
    msg: '',
    svcs: []
  },
  methods:{
    get() {
      axios({ method: "GET", "url": "/api/deployments" }).then(result => {
        this.svcs = result.data;
      }, error => {
        console.error(error);
      });
    },
    change(state, ns, name) {
      axios({
        method: "PUT",
        url: `/api/deployments/${ns}/${name}/${state}`
      }).then(result => {
        this.msg = `${name} (${ns}): going ${state}`;
        this.get();
        window.setTimeout(() => {
          this.msg = ''
        }, 5000);
      }, error => {
        console.error(error);
      });
    }
  },
  beforeMount(){
    this.get()
  }
})
