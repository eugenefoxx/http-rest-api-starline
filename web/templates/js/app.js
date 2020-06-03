console.log('work')
const app = new Vue({
    el: '#app',
    data: {
        errors: [],
        material: null,
        qty: null,
        comment: null
    },
    methods: {
        shipmentBySAP: function (e) {
            /*if (this.material && this.qty) {
                return true;
            }

            this.errors = [];

            if (!this.name) {
                this.errors.push('Требуется указать material.');
            }
            if (!this.age) {
                this.errors.push('Требуется указать qty.');
            }*/
            console.log(e)
            //e.preventDefault();

        }
    }
})