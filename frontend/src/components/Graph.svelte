<script lang="ts">
  import { onMount } from "svelte";
  export let graphData: { x: number[], y: number[], type: string };

  let plotDiv;
  onMount(() => {
    let Plot = new Plotly.newPlot(
      plotDiv,
      [graphData],
      {
        title: false,
        height: 200,
        margin: { r: 60, l: 0, t: 20, b: 20, pad: 5 },
        colorway: ['#c9c9c9'],
        yaxis: { side: "right", ticksuffix: "ms", color: "#949494" },
        xaxis: { color: "white" }
      },
      { showSendToCloud: false, staticPlot: true }
    );
  });
</script>

<div id="plotly">
  <div id="plotDiv" bind:this={plotDiv}><!-- Plotly chart will be drawn inside this DIV --></div>
</div>
