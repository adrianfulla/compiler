<template>
    <div class="textarea-wrapper ">
      <textarea v-model="inputData" @input="handleInput" placeholder="Escribe aquí..."></textarea>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      sendSignal: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        inputData: ''
      };
    },
    methods: {
      handleInput() {
        if (this.sendSignal) {
          this.$emit('sendData', this.inputData);
        }
      }
    },
    watch: {
      sendSignal(newVal) {
        if (newVal) {
          this.$emit('sendData', this.inputData);
        }
      }
    }
  }
  </script>
  
  <style scoped>
  .textarea-wrapper {
    width: 100%;
    height: 100%;
  }
  
  textarea {
    width: 100%;
    height: 100%;
    border: dark;
    padding: 10px; /* Añade un poco de padding para que el texto no esté justo en el borde */
    box-sizing: border-box; /* Asegura que el padding y el borde estén incluidos en la dimensión */
    resize: none; /* Evita que el usuario redimensione el textarea manualmente */
  }
  </style>
  