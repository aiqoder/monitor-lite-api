<template>
    <div ref="yamlEditor" style="height: 100%;"></div>
</template>

<script lang="ts" setup>
import { onMounted, ref, unref } from 'vue';
import { basicSetup } from "codemirror";
import { EditorState } from "@codemirror/state"
import { EditorView } from "@codemirror/view"
import { yaml } from "@codemirror/lang-yaml"
import { oneDark } from '@codemirror/theme-one-dark'
import { watch } from 'vue';

const yamlEditor = ref()
const doc = defineModel("value", { default: () => "" })

let myTheme = EditorView.theme({
    "&": {
        color: "white",
        backgroundColor: "#034"
    },
    ".cm-content": {
        caretColor: "#0e9"
    },
    "&.cm-focused .cm-cursor": {
        borderLeftColor: "#0e9"
    },
    "&.cm-focused .cm-selectionBackground, ::selection": {
        backgroundColor: "#074"
    },
    ".cm-gutters": {
        backgroundColor: "#045",
        color: "#ddd",
        border: "none"
    }
}, { dark: true })

let editor: EditorView;
let timer;
async function initRuleEditorEntity(dom: HTMLDivElement) {
    const startState = EditorState.create({
        doc: '',
        extensions: [
            basicSetup,
            oneDark,
            yaml(),
            myTheme,
            EditorView.updateListener.of((v) => {
                const rule = v.state.doc.toString()
                if (!rule) return
                doc.value = rule
            }),
        ],

    })

    if (dom) {
        editor = new EditorView({
            state: startState,
            parent: dom,
        });
    }
}

onMounted(() => {
    initRuleEditorEntity(unref(yamlEditor))
})

let init = 0;
watch(doc, (val) => {
    if (!val) return
    doc.value = val
    if(init > 0) return
    init++
    editor.dispatch({
        changes: {
            from: 0,
            to: editor.state.doc.length,
            insert: val
        }
    })
})
</script>
<style lang="scss" scoped></style>