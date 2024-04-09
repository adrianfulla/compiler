<template>
    <div v-if="treeImage" class="tree">
      <img :src="treeImage" alt="Imagen de arbol de expresion generada">
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
        treeImage: null,
    }
        
    },
    async created() {
      try {
        const firstResponse = await fetch('http://localhost:8080/automata/arbol',{
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({regex: this.regex})
        })
        if (firstResponse.ok) {
            const treeData = await firstResponse.json()
            this.fetchImage(treeData)
        }
        else {
            console.error('Error al obtener los datos: ', response.statusText)
        }
      } catch (error) {
        console.error('Error al hacer la solicitud: ', error)
      }
    },
    methods:{
      async fetchImage(treeData) {
        try {
          const response = await fetch('http://localhost:5000/arbol',{
            method: 'POST',
                  headers: {
                    'Content-Type': 'application/json'
                  },
                  body: JSON.stringify(treeData)
          }); // Asume que esta es la URL de tu API
          if (response.ok) {
            const blob = await response.blob();
            this.treeImage = URL.createObjectURL(blob);
          } else {
            console.error('Error al obtener la imagen:', response.statusText);
          }
        } catch (error) {
          console.error('Error en la solicitud:', error);
        }
    }
    }
  };
  </script>
  