<template>
    <h1>Generador de automata</h1>
    <div class="input-group mb-3">
        <input type="text" class="form-control" placeholder="Ingrese la expresion regular" v-model="regex">
        <button class="btn btn-primary" @click="generateTree">Generar</button>
    </div>
    <div class="card" v-if="mostrar_automata">
      <div class="container mt-4" >
        <div class="row">
          <div class="col-md-4" >
            <div class="card">
            <div class="card-header">
              Árbol de Nodos
            </div>
            <div class="card-body">
              <NodeTree :regex="regex"/>
            </div>
          </div>
        
      </div>
      <div class="col-md-8" >
        <div class="card">
          <div class="card-header">
            Automata Determinista
          </div>
          <div class="card-body">
            <DirAFD :regex="regex"/>
          </div>
        </div>
        </div>  
      </div>
      </div>
    </div>
  </template>
  
  <script>
  import NodeTree from '../components/Automata/NodeTree.vue'; // Asegúrate de que la ruta del import sea correcta
  import DirAFD from '../components/Automata/DirAFD.vue'
  export default {
    name: 'TreeView',
    components: {
      NodeTree,
      DirAFD
    },
    data() {
      return {
        regex: '',
        mostrar_automata: false 
      };
    },
    methods:{
        generateTree(){
                if (this.regex.length > 0) {
                this.mostrar_automata = true;
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
  