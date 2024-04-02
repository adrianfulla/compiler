<template>
    <h1>Generador de automata</h1>
    <div class="input-group mb-3">
        <input type="text" class="form-control" placeholder="Ingrese la expresion regular" v-model="regex">
        <button class="btn btn-primary" @click="generateTree">Generar</button>
    </div>
    
    <div class="container mt-4" v-if="mostrarArbol">
      <div class="card">
        <div class="card-header">
          Árbol de Nodos
        </div>
        <div class="card-body">
          <NodeTree :regex="regex"/>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import NodeTree from '../components/Automata/NodeTree.vue'; // Asegúrate de que la ruta del import sea correcta
  
  export default {
    name: 'TreeView',
    components: {
      NodeTree
    },
    data() {
      return {
        regex: '',
        mostrarArbol: false 
      };
    },
    methods:{
        generateTree(){
                if (this.regex.length > 0) {
                this.mostrarArbol = true;
            }
        }
    },
    created() {
    // Aquí podrías hacer la llamada a la API para obtener los datos del árbol si es necesario
    // o utilizar una tienda de Vuex/Pinia para obtener los datos del árbol
    }
}
  </script>
  
  <style>
  /* Aquí puedes añadir estilos específicos para esta vista */
  </style>
  