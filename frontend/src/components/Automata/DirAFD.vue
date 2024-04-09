<template>
    <div v-if="!mostrar_afd">
        Cargando...
    </div>
  <div id="dir_afd" v-else>
      <img :src="dir_afd_image" alt="Imagen de automata directo generado">
  </div>
</template>

<script>
import GraphNode from './GraphNode.vue'
export default {
    name: 'dir_afd',
    props:['regex'],
    components: {
      GraphNode
    },
  data() {
    return {
        dir_afd_image: null,
        mostrar_afd: false,
        dir_afd: null,
    };
  },
    async created() {
      try {
        const firstResponse = await fetch('http://localhost:8080/automata/afd',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({regex: this.regex})
        })
        if (firstResponse.ok) {
            this.dir_afd = await firstResponse.json()
            this.fetchImage(this.dir_afd)
            this.emitValue()
        }
      } catch (error) {
        console.error('Error al hacer la solicitud: ', error)
      }
    },
    methods: {
      async fetchImage(dir_afd) {
        try {
          const response = await fetch('http://localhost:5000/afd',{
            method: 'POST',
                  headers: {
                    'Content-Type': 'application/json'
                  },
                  body: JSON.stringify(dir_afd)
          }); // Asume que esta es la URL de tu API
          if (response.ok) {
            const blob = await response.blob();
            this.dir_afd_image = URL.createObjectURL(blob);
            this.mostrar_afd = true
          } else {
            console.error('Error al obtener la imagen:', response.statusText);
          }
        } catch (error) {
          console.error('Error en la solicitud:', error);
        }
    },
    emitValue(){
      this.$emit('afd', this.dir_afd)
    }
  },
  };
</script> 
  