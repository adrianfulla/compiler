<template>
    <div v-if="treeData" class="tree">
      <Node :data="treeData.raiz" />
    </div>
    <div v-else>
        Cargando...
    </div>
  </template>
  
  <script>
  import Node from './Node.vue';
  
  export default {
    name: 'Tree',
    props:['regex'],
    components: {
      Node
    },
    data() {return {
        treeData: null,
    }
        
    },
    async created() {
      try {
        const response = await fetch('http://localhost:8080/automata/arbol',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({regex: this.regex})
        })
        if (response.ok) {
            this.treeData = await response.json()
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
  