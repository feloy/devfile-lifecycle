import mermaid from 'https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.esm.min.mjs';

export function onInit() {
    mermaid.initialize({startOnLoad: false});

    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("devfile.wasm"), go.importObject).then((result) => {
        go.run(result.instance);                
    });
}
   
export async function onButtonClick() {
    const source = document.querySelector("#input").value;
    const element = document.querySelector("#mermaid");
    const result = getFlowChart(source);
    console.log(result);
    const { svg } = await mermaid.render('rendered', result);
    console.log(svg);
    element.innerHTML = svg;
}
