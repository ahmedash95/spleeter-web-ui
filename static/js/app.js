var app = new Vue({
    el: '#app',
    data: {
        file: null,
        progress: 0,
        processing: false,
        results: {
            vocals: null,
            accompaniment: null,
        }
    },
    computed: {
        uploadable(){
            return !this.processing && this.file === null;
        }
    },
    methods: {
        fileUpload(event) {
            this.file = this.$refs.file.files[0];
        },
        submit(){
            let formData = new FormData();
            formData.append('file', this.file);
            this.processing = true;
            axios.post('/upload',formData,{
                headers: {
                    'Content-Type': 'multipart/form-data'
                },
                onUploadProgress: function( progressEvent ) {
                    this.progress = parseInt( Math.round( ( progressEvent.loaded * 100 ) / progressEvent.total ) );
                }.bind(this)
            }).then((resp) => {
                this.processing = false;
                this.results.accompaniment = resp.data.accompaniment
                this.results.vocals = resp.data.vocals
            }).finally(() => {
                this.processing = false;
            });
        }
    }
})