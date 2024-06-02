<template>
    <h1>Generador de automata</h1>
    <div class="input-group mb-3">
        <input type="text" class="form-control" placeholder="Ingrese la expresion regular" v-model="regex">
        <button class="btn btn-primary" @click="generateTree">Generar</button>
    </div>
    <div class="card" v-if="mostrar_automata">
      <div class="container mt-4 mb-4" >
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
            <DirAFD :regex="regex" @afd="handleAfdReturn"/>
          </div>
        </div>
        </div>  
      </div>
      <div class="card">
        <div class="card-header">
          Ingrese la cadena a simular
        </div>
        <div class="card-body">
          <TextBlock :sendSignal="sendSimulateString" @sendData="receiveSimulateString"/>
          <button class="btn btn-success float-center" @click="triggerSendSimulateString">Simular</button>
        </div>
        <div class="card-footer" v-if="simulationReturn">
            <div class="alert alert-success" v-if="simulationResult">
                La cadena fue aceptada por el automata
            </div>
            <div class="alert alert-warning" v-else>
                La cadena no fue aceptada por el automata
            </div>
        </div>
      </div>
      </div>
    </div>
  </template>
  
  <script>
  import NodeTree from '../components/Automata/NodeTree.vue'; // Asegúrate de que la ruta del import sea correcta
  import DirAFD from '../components/Automata/DirAFD.vue'
  import TextBlock from '../components/LexAnalyzer/TextBlock.vue';
  
  export default {
    name: 'TreeView',
    components: {
      NodeTree,
      DirAFD,
      TextBlock
    },
    data() {
      return {
        regex: '',
        mostrar_automata: false,
        afd: null,
        sendSimulateString: false,
        simulationReturn: false,
        simulationResult: null
      };
    },
    methods:{
        generateTree(){
                if (this.regex.length > 0) {
                this.mostrar_automata = true;
            }
        },
        handleAfdReturn(afd){
          console.log(afd)
          this.afd = afd
        },
        triggerSendSimulateString() {
          this.sendSimulateString = true;
          // console.log("ACA")
    
          // Restablecer la señal después de activarla
          this.$nextTick(() => {
            this.sendSimulateString = false;
          })
        },
        receiveSimulateString(data){
          this.simulateString(data)
        },

        async simulateString(data){
          try { 
              const response = await fetch('http://localhost:8080/automata/afd/simulate/',{
                method: 'POST',
                      headers: {
                        'Content-Type': 'application/json'
                      },
                      body: JSON.stringify({Regex:this.regex, Simulate: data})
              }); 
              if (response.ok) {
                this.simulationReturn = true
                this.simulationResult = await response.json()
              } else {
                console.error('Error, hubo un error al simular');
              }
            } catch (error) {
              console.error('Error en la solicitud:', error);
            }
        },
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
  