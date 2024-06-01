<template>
    <div>
        <h3>Analizador Lexico</h3>
    </div>
    <div>
        <div class="card mb-2" id="yalex-text">
            <div class="card-header">
                Configuraciones YaLex
            </div>
            <div class="card-body mb-4" style="height: 500px;">
                <div class="row">
                    <div class="col-6" style="height:450px">
                    <h5>YaLex</h5>
                    <TextBlock :sendSignal="sendSignalLex" @sendData="receiveDataYalex" id="yalex"/>
                </div>
                <div class="col-6" style="height: 450px;">
                    <h5>YaPar</h5>
                    <TextBlock :sendSignal="sendSignalLex" @sendData="receiveDataYapar" id="yapar"/>
                </div>
            </div>
        </div>
        <div class="card-footer">
            <button class="btn btn-success float-center" @click="triggerSendFiles">Validar Yalex y YaPar</button>
        </div>
        </div>
        <div class="card mb-2" v-if="yaparExitoso">
          <div class="card-header">
            SLR - Grafico
            <div v-if="!slr_image">
              ...cargando
            </div>
            <div v-else>
              <img :src="slr_image" alt="Imagen de slr generado" style="max-width: 100%; max-height: 100%;">
            </div>
          </div>
          <div class="card-body">
            SLR - Tabla
            <SLRTable :slr_table="slr_table"></SLRTable>
            LR(1) - Tabla de parseo
            <LR1Table :lr1Table="lr1_table"></LR1Table>
          </div>
          <div class="card-footer">
            Ingrese lo que se desea parsear
            <TextBlock :sendSignal="sendParsingSignal" @sendData="receiveDataParsing" id="parsing"/>
            <button class="btn btn-success float-center" @click="triggerSendParsing">Parsear Cadena</button>
            <div class="alert alert-info text-center ">
              {{ parsingResponse }}
            </div>
          </div>
        </div>
    </div>
</template>
<script>
import TextBlock from '../components/LexAnalyzer/TextBlock.vue';
import SLRTable from '../components/LexAnalyzer/SLRTable.vue';
import LR1Table from '../components/LexAnalyzer/LR1Table.vue'

export default {
  components: {
    TextBlock,
    SLRTable,
    LR1Table
  },
  data() {
    return {
      sendSignalLex: false,
      sendParsingSignal: false,
      yaparExitoso: false,
      yalex1: null,
      yapar: null,
      parsingString: null,
      tokens: [],
      slr_table: null,
      slr_image: null,
      lr1_table: null,
      parsingResponse: null,
    };
  },
  methods: {
    triggerSendFiles() {
      this.sendSignalLex = true;
      // Restablecer la señal después de activarla
      this.$nextTick(() => {
        this.sendSignalLex = false;
      });
    },
    triggerSendParsing() {
      this.sendParsingSignal = true;
      // Restablecer la señal después de activarla
      this.$nextTick(() => {
        this.sendParsingSignal = false;
      });
    },
    receiveDataYalex(data) {
      // console.log('Datos recibidos de Yalex 1:', data);
      this.yalex1 = data;
    },
    receiveDataYapar(data) {
      // console.log('Datos recibidos de Yalex 2:', data);
      this.yapar = data;
      this.validate();
    },
    receiveDataParsing(data) {
      console.log('Datos recibidos de parsing:', data);
      this.parsingString = data;
      this.sendParsing();
    },
    validate(){
      this.fetchSLR()
    },
    async fetchSLR(){
      try {
        const response = await fetch(`http://localhost:8080/parser/slr/`,{
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({yalex: this.yalex1, yapar: this.yapar})
        });
        if (response.ok) {
          this.yaparExitoso = true;
          // console.log(response, response.json())
          const responseJson = await response.json()
          console.log(responseJson)
          this.slr_table = responseJson.print_slr
          this.lr1_table = responseJson.lr1_table
          // console.log(this.lr1_table, response.lr1_table)
          this.fetchSLRImage(await responseJson.slr)
        } else {
          console.error(`Error, Yapar invalido`);
        }
      } catch (error) {
        console.error('Error en la solicitud:', error);
      }
    },

    async sendParsing(){
      try {
        const response = await fetch(`http://localhost:8080/parser/lr1/`,{
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({yalex: this.yalex1, yapar: this.yapar, parsing: this.parsingString})
        });
        if (response.ok) {
          // console.log(response, response.json())
          const responseJson = await response.json()
          console.log(responseJson)
          this.parsingResponse = responseJson
        } else {
          
          this.parsingResponse = await response.json()
          console.error(`Error, parsing string invalido`, this.parsingResponse );
        }
      } catch (error) {
        console.error('Error en la solicitud:', error);
      }
    },

    async fetchSLRImage(data){
      try {
          const response = await fetch('http://localhost:5000/lr0',{
            method: 'POST',
                  headers: {
                    'Content-Type': 'application/json'
                  },
                  body: JSON.stringify(data)
          }); // Asume que esta es la URL de tu API
          if (response.ok) {
            const blob = await response.blob();
            this.slr_image = URL.createObjectURL(blob);
            this.mostrar_slr = true
          } else {
            console.error('Error al obtener la imagen:', response.statusText);
          }
        } catch (error) {
          console.error('Error en la solicitud:', error);
        }
    }
  }
}
</script>
