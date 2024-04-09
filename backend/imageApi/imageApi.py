from flask import Flask, request, send_file
from graphviz import Digraph
from datetime import datetime

app = Flask(__name__)

class Nodo:
    def __init__(self, valor, izquierdo=None, derecho=None, leafNum=None, nullability=False, firstpos=None, lastpos=None, followpos=None):
        self.valor = valor
        self.izquierdo = izquierdo
        self.derecho = derecho
        self.leaf = leafNum
        self.nullability = nullability
        self.firstpos = firstpos
        self.lastpos = lastpos
        self.followpos = followpos or []

    def repr(self):
        return f"{self.valor}\n"

    @staticmethod
    def desde_dict(d):
        valor = d.get('valor')
        leaf = d.get('leaf', None)
        firstpos = d.get('firstpos', [])
        lastpos = d.get('lastpos', [])
        followpos = d.get('followpos', [])
        izquierdo = Nodo.desde_dict(d['izquierdo']) if 'izquierdo' in d else None
        derecho = Nodo.desde_dict(d['derecho']) if 'derecho' in d else None
        return Nodo(valor, izquierdo, derecho, leaf, firstpos=firstpos, lastpos=lastpos, followpos=followpos)

def visualizar_arbol(raiz):
    def agregar_nodos_edges(raiz, dot=None):
        if dot is None:
            dot = Digraph()
            dot.node(name=str(raiz), label=raiz.repr())

        if raiz.izquierdo:
            dot.node(name=str(raiz.izquierdo), label=raiz.izquierdo.repr())
            dot.edge(str(raiz), str(raiz.izquierdo))
            agregar_nodos_edges(raiz.izquierdo, dot)

        if raiz.derecho:
            dot.node(name=str(raiz.derecho), label=raiz.derecho.repr())
            dot.edge(str(raiz), str(raiz.derecho))
            agregar_nodos_edges(raiz.derecho, dot)

        return dot

    return agregar_nodos_edges(raiz)

class Dstate:
    def __init__(self, nombre, aceptacion=False):
        self.nombre = nombre
        self.aceptacion = aceptacion
        self.transiciones = {}

    def add_transicion(self, simbolo, estado):
        self.transiciones[simbolo] = estado

    def __repr__(self):
        return f"Dstate({self.nombre}, Accept={self.aceptacion})"


class Afd:
    def __init__(self, estados, alfabeto, estado_inicial, estados_finales, transiciones):
        self.estados = estados
        self.alfabeto = alfabeto
        self.estado_inicial = estado_inicial
        self.estados_finales = estados_finales
        self.transiciones = transiciones

    def visualizar_afd(self):
        dot = Digraph(comment='AFD')
        dot.attr(rankdir='LR')
        dot.node('start', '', shape='point', style='invisible')

        # Agrega los nodos al gráfico
        for estado in self.estados:
            if estado in self.estados_finales:
                dot.node(estado, shape='doublecircle')
            else:
                dot.node(estado)

        # Agrega un nodo inicial invisible para apuntar al estado inicial
        dot.edge('start', self.estado_inicial, style='bold')

        # Agrega las transiciones al gráfico
        for origen, destinos in self.transiciones.items():
            for simbolo, destino in destinos.items():
                dot.edge(origen, destino, label=simbolo)

        return dot





@app.route('/')
def hello_world():
    return 'Hello, World!'

@app.route('/arbol', methods=['POST'])
def arbol():
    data = request.get_json()
    raiz_dict = data.get('raiz')
    if raiz_dict:
        raiz = Nodo.desde_dict(raiz_dict)
        dot = visualizar_arbol(raiz)
        now = datetime.now()
        dot.render('./images/arbol_expresion'+now.strftime("%H:%M:%S"), view=False, format='jpg')
        return send_file('./images/arbol_expresion'+now.strftime("%H:%M:%S")+'.jpg', mimetype='image/jpg')
    return {'recibido':data}

@app.route('/afd', methods=['POST'])
def afd():
    data = request.get_json()
    afd = Afd(data['estados'], data['alfabeto'], data['estado_inicial'],data['estados_finales'], data['transiciones'])
    if afd:

        dot = afd.visualizar_afd()
        now = datetime.now()
        dot.render('./images/afd'+now.strftime("%H:%M:%S"), view=False, format='jpg')
        return send_file('./images/afd'+now.strftime("%H:%M:%S")+'.jpg', mimetype='image/jpg')
    return {'recibido':data}

if __name__ == '__main__':
    app.run(host='localhost', port=5000)

