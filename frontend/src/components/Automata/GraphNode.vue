<template>
    <div class="tooltip-wrapper">
      <v-network-graph
        ref="graph"
        :nodes="nodes"
        :edges="edges"
        :configs="options"
        :event-handlers="eventHandlers"
      />
      <!-- Tooltip -->
      <div
        ref="tooltip"
        class="tooltip"
        :style="{ left: tooltipPos.left, top: tooltipPos.top, opacity: tooltipOpacity }"
      >
        <div v-if="targetNodeType === 'node'">
          {{ tooltipContent.name }}
          <br />
          <span v-if="tooltipContent.isFinal">Estado de aceptación</span>
        </div>
        <div v-if="targetNodeType === 'edge'">
          {{ tooltipContent.label }}
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref, computed, watch } from 'vue'
  import { VNetworkGraph, defineConfigs } from 'v-network-graph'
  
  export default {
    components: {
      VNetworkGraph
    },
    props: {
      estados: Array,
      alfabeto: Array,
      estado_inicial: String,
      estados_finales: Array,
      transiciones: Object
    },
    setup(props) {
      const graph = ref(null)
      const tooltip = ref(null)
      const tooltipOpacity = ref(0)
      const tooltipPos = ref({ left: '0px', top: '0px' })
      const targetNodeId = ref(null)
      const targetEdgeId = ref(null)
      const targetNodeType = ref('')
    
      const nodes = computed(() => {
        return props.estados.map(estado => ({
          id: estado,
          label: estado,
          name: estado,
          class: props.estados_finales.includes(estado) ? 'final' : ''
        }));
      });
  
      const edges = computed(() => {
        const edges = [];
        const stateIndexMap = {};
        console.log(nodes)
        nodes.value.forEach((node, index) => {
            stateIndexMap[node.id] = index
        })
        for (const [origen, destinos] of Object.entries(props.transiciones)) {
          for (const [simbolo, destino] of Object.entries(destinos)) {
            edges.push({
              source: stateIndexMap[origen],
              target: stateIndexMap[destino],
              label: simbolo,
            });
          }
        }
        return edges;
      });
  
      const tooltipContent = computed(() => {
        if (targetNodeType.value === 'node') {
          const node = nodes.value.find(n => n.id === targetNodeId.value)
          return {
            name: node.label,
            isFinal: props.estados_finales.includes(node.id)
          };
        } else if (targetNodeType.value === 'edge') {
          const edge = edges.value.find(e => e.id === targetEdgeId.value)
          return {
            label: edge.label
          };
        }
        return {};
      });
  
      const eventHandlers = {
        'node:pointerover': ({ node }) => {
          targetNodeId.value = node
          targetNodeType.value = 'node'
          tooltipOpacity.value = 1
        },
        'node:pointerout': () => {
          tooltipOpacity.value = 0
        },
        'edge:pointerover': ({ edge }) => {
          targetEdgeId.value = edge
          targetNodeType.value = 'edge'
          tooltipOpacity.value = 1
        },
        'edge:pointerout': () => {
          tooltipOpacity.value = 0
        }
      };
  
      watch([targetNodeId, tooltipOpacity], () => {
        if (!graph.value || !tooltip.value) return
  
        const nodePos = graph.value.getNodeRect(targetNodeId.value)
        if (nodePos) {
          const domPoint = graph.value.translateFromSvgToDomCoordinates({
            x: nodePos.x + nodePos.width / 2,
            y: nodePos.y
          })
          tooltipPos.value = {
            left: domPoint.x - tooltip.value.offsetWidth / 2 + 'px',
            top: domPoint.y - tooltip.value.offsetHeight - 10 + 'px'
          }
        }
      }, { immediate: true });

      const options = computed(() => {
        return defineConfigs({
            node: {
            class: (node) => node.class,
            size: 10
          },
          edges: {
            label: {
              position: 'center',
              text: (edge) => edge.label
            },
            markerEnd: { type: 'arrow' },
            summarized:{
                label:{
                    fontSize: 10,
                    color: '#4466cc'
                },
                shape: {
                    type: 'circle',
                    radius: 6,
                    color: '#ffffff',
                    strokeWidth: 1,
                    strokeColor: '#4466cc'
                }
            },
            stroke: {
              width: 5,
              color: '#4466cc'
            }
            }
        })
        })
  
      return { graph, tooltip, nodes, edges, options, eventHandlers, tooltipPos, tooltipOpacity, tooltipContent, targetNodeType };
    }
}
  </script>
  
  <style scoped>
  .tooltip-wrapper {
    position: relative;
  }
  .tooltip {
    position: absolute;
    padding: 10px;
    background-color: #fff;
    border: 1px solid #ccc;
    box-shadow: 0px 0px 5px rgba(0,0,0,0.2);
    border-radius: 5px;
    opacity: 0;
    transition: opacity 0.3s;
    pointer-events: none;
  }
  </style>
  