<template>
    <div class="node">
      <div>{{ data.valor }}</div>
      <div class="children">
        <Node v-if="data.izquierdo" :data="data.izquierdo" />
        <Node v-if="data.derecho" :data="data.derecho" />
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'Node',
    props: {
      data: Object
    }
  };
  </script>
  
  <style>
  .node {
    margin: 10px;
    padding: 10px;
    border: 1px solid #ccc;
  }
  
  .children {
    margin-left: 20px;
  }
  </style>
  