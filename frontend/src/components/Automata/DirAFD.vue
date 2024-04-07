<template>
    <div v-if="!mostrar_afd">
        Cargando...
    </div>
  <div id="dir_afd" v-else>
    <GraphNode :estados='dir_afd.estados' :alfabeto='dir_afd.alfabeto' :estado_inicial="dir_afd.estado_inicial"  :estados_finales='dir_afd.estados_finales'  :transiciones='dir_afd.transiciones' ></GraphNode>
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
        dir_afd: null,
        mostrar_afd: false,
    };
  },
    async created() {
      try {
        const response = await fetch('http://localhost:8080/automata/afd',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({regex: this.regex})
        })
        if (response.ok) {
            this.dir_afd = await response.json()
            console.log(this.dir_afd)
            this.mostrar_afd = true
        }
        else {
            console.error('Error al obtener los datos: ', response.statusText)
        }
      } catch (error) {
        console.error('Error al hacer la solicitud: ', error)
      }
    }
  };
</script> 
  