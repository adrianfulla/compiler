<template>
    <div>
        <h3>
            Analizador Lexico
        </h3>
    </div>
    <div>
        <div class="card mb-2" id="yalex-text">
            <div class="card-header">
                YaLex
            </div>
            <div class="card-body" style="width: 100%; height: 500px;" id="yalex-text">
                <TextBlock :sendSignal="sendSignalYalex" @sendData="receiveDataYalex" id="yalex-text"/>
            </div>
            <div class="card-footer ">
                <button class="btn btn-success float-center" @click="triggerSendYalex">Validar Yalex</button>
            </div>
        </div>
        <div class="card" id="analize-text" v-if="yalexExitoso">
            <div class="card-header">
                Analizador
            </div>
            <div class="card-body">
                <TextBlock :sendSignal="sendSignalLex" @sendData="receiveDataLex"/>
            </div>
            <div class="card-footer ">
                <button class="btn btn-success float-center" @click="triggerSendLex">Analizar</button>
                <div v-if="tokens.length > 0">
                    <h3>Resultados del análisis:</h3>
                    <table class="table">
                        <thead>
                            <tr>
                                <th>Token</th>
                                <th>Valor</th>
                                <th>Inicio</th>
                                <th>Fin</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="token in tokens" :key="token.start">
                                <td>{{ token.token }}</td>
                                <td>{{ token.value }}</td>
                                <td>{{ token.start }}</td>
                                <td>{{ token.end }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
  </template>
  
  <script>
import TextBlock from '../components/LexAnalyzer/TextBlock.vue';
  
  export default {
    components: {
      TextBlock
    },
    data() {
      return {
        sendSignalYalex: false,
        sendSignalLex: false,
        yalexExitoso: false,
        yalex: null,
        tokens: [],
      };
    },
    methods: {
     triggerSendYalex() {
        this.sendSignalYalex = true;
        // console.log("ACA")
  
        // Restablecer la señal después de activarla
        this.$nextTick(() => {
          this.sendSignalYalex = false;
        });
      },
      triggerSendLex() {
        this.sendSignalLex = true;
  
        // Restablecer la señal después de activarla
        this.$nextTick(() => {
          this.sendSignalLex = false;
        });
      },
      receiveDataYalex(data) {
        // console.log('Yalex Datos recibidos:', data);
        this.validateYalex(data)
      },
      receiveDataLex(data) {
        console.log('Lex Datos recibidos:', data)
        this.parseLex(data)
      },

      async validateYalex(data){
        try {
          const response = await fetch('http://localhost:8080/lexer/create',{
            method: 'POST',
                  headers: {
                    'Content-Type': 'application/json'
                  },
                  body: JSON.stringify({yalex: data})
          }); // Asume que esta es la URL de tu API
          if (response.ok) {
            this.yalexExitoso = true
            this.yalex = data
          } else {
            console.error('Error, YaLex no valido, hubo un error al crear el Scanner');
          }
        } catch (error) {
          console.error('Error en la solicitud:', error);
        }
    },

    async parseLex(data){
        try {
          const response = await fetch('http://localhost:8080/lexer/parselex',{
            method: 'POST',
                  headers: {
                    'Content-Type': 'application/json'
                  },
                  body: JSON.stringify({yalex: this.yalex, parsing: data + " "})
          }); // Asume que esta es la URL de tu API
          if (response.ok) {
            const result = await response.json()
            this.tokens = result.reverse()

          } else {
            console.error('Error, expresion no pudo ser analizada');
            this.tokens = []
          }
        } catch (error) {
          console.error('Error en la solicitud:', error);
          this.tokens = []
        }
    },

    }
  }
  </script>
  