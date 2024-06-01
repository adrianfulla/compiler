<template>
    <div class="card">
      <div class="card-header">
        Detalles del SLR
        <button class="btn btn-link" @click="toggleCollapse">{{ collapsed ? 'Expandir' : 'Colapsar' }}</button>
      </div>
      <div class="card-body" v-if="!collapsed">
        <table class="table table-striped">
          <thead>
            <tr>
              <th>Estado ID</th>
              <th>Items</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="state in slr_table" :key="state.state_id">
              <td>{{ state.state_id }}</td>
              <td>
                <ul>
                  <li v-for="item in state.items" :key="item.description">
                    {{ item.description }}
                  </li>
                </ul>
              </td>
              <td>
                <ul>
                  <li v-for="action in state.actions" :key="action">
                    {{ action }}
                  </li>
                </ul>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'SLRTable',
    props: {
      slr_table: Array
    },
    data() {
      return {
        collapsed: true // Comienza colapsado
      }
    },
    methods: {
      toggleCollapse() {
        this.collapsed = !this.collapsed;
      }
    }
  }
  </script>
  