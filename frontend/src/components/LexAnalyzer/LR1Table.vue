<template>
    <div class="card">
      <div class="card-header">
        Detalles del LR(1)
        <button class="btn btn-link" @click="toggleCollapse">{{ collapsed ? 'Expandir' : 'Colapsar' }}</button>
      </div>
      <div class="card-body" v-if="!collapsed">
        <table class="table table-striped">
          <thead>
            <tr>
              <th>Estado ID</th>
              <th>Ítems</th>
              <th>GOTO</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="state in lr1Table.states" :key="state.id">
              <td>{{ state.id }}</td>
              <td>
                <ul>
                  <li v-for="item in state.items" :key="item.production.head + item.position">
                    {{ item.production.head }} → {{ item.production.body[0].join(' ') }} ({{ item.lookaheads.join(', ') }})
                  </li>
                </ul>
              </td>
              <td>
                <ul v-if="lr1Table.gotos && lr1Table.gotos[state.id]">
                  <li v-for="(dest, symbol) in lr1Table.gotos[state.id]" :key="symbol">
                    {{ symbol }}: Estado {{ dest }}
                  </li>
                </ul>
                <p v-else>No transitions</p>
              </td>
              <td>
                <ul v-if="lr1Table.actions && lr1Table.actions[state.id]">
                  <li v-for="(action, symbol) in lr1Table.actions[state.id]" :key="symbol">
                    {{ symbol }}: {{ action }}
                  </li>
                </ul>
                <p v-else>No actions</p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
</template>

<script>
export default {
  name: 'LR1Table',
  props: {
    lr1Table: Object
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
