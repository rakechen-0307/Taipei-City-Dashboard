<script setup>
import { ref, computed } from "vue";
import { useMapStore } from "../../store/mapStore";

const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);
// const mapStore = useMapStore();

const targetArea = ref(null);
// const AreaColor = ref(props.chart_config.color[0]);
const mousePosition = ref({ x: null, y: null });
const refMousePosition = ref({ x: null, y: null });
const selectedIndex = ref(null);
const mapStore = useMapStore();

// Parse Venn Diagram Data
const VennData = computed(() => {
	let readData = {};
	let count = {};
	let sum = 0;

	let i = 0;
	while (props.series[i].name !== null) {
		i++;
	}
	i++;

	let len = props.series[i].data.length;
	if (len === 2) {
		for (let i = 0; i < 3; i++) {
			count[i.toString()] = 0;
		}
	} else if (len === 3) {
		for (let i = 0; i < 7; i++) {
			count[i.toString()] = 0;
		}
	}

	while (i < props.series.length) {
		sum++;
		const item = props.series[i];
		const dat = item.data;
		readData[item.name] = dat;
		if (len === 2) {
			if ((dat[0] === "1") & (dat[1] === "1")) {
				count["2"]++;
			} else if ((dat[0] === "1") & (dat[1] === "0")) {
				count["0"]++;
			} else if ((dat[0] === "0") & (dat[1] === "1")) {
				count["1"]++;
			}
		} else if (len === 3) {
			if ((dat[0] === "1") & (dat[1] === "1") & (dat[2] === "1")) {
				count["6"]++;
			} else if ((dat[0] === "1") & (dat[1] === "0") & (dat[2] === "0")) {
				count["0"]++;
			} else if ((dat[0] === "0") & (dat[1] === "1") & (dat[2] === "0")) {
				count["1"]++;
			} else if ((dat[0] === "0") & (dat[1] === "0") & (dat[2] === "1")) {
				count["2"]++;
			} else if ((dat[0] === "1") & (dat[1] === "1") & (dat[2] === "0")) {
				count["3"]++;
			} else if ((dat[0] === "1") & (dat[1] === "0") & (dat[2] === "1")) {
				count["4"]++;
			} else if ((dat[0] === "0") & (dat[1] === "1") & (dat[2] === "1")) {
				count["5"]++;
			}
		}
		i++;
	}

	const output = { sum: sum, len: len, readData: readData, count: count };
	return output;
});

const tooltipPosition = computed(() => {
	return {
		left: `${mousePosition.value.x - 10}px`,
		top: `${mousePosition.value.y - 54}px`,
	};
});
function toggleActive(e) {
	targetArea.value = e.target.dataset.name;
}
function toggleActiveToNull() {
	targetArea.value = null;
}
function updateMouseLocation(e) {
	const container = document.getElementsByClassName("svg");
	const rect = container[0].getBoundingClientRect();
	refMousePosition.value.x =
		((e.pageX - rect.left) / (rect.right - rect.left)) * 550;
	refMousePosition.value.y =
		((e.pageY - rect.top) / (rect.bottom - rect.top)) * 550;
	mousePosition.value.x = e.pageX;
	mousePosition.value.y = e.pageY;
}

function handleDataSelection(index) {
	if (!props.chart_config.map_filter) {
		return;
	}
	if (index !== selectedIndex.value) {
		mapStore.addLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`,
			props.chart_config.map_filter[2],
			props.chart_config.map_filter[3][index]
		);
		selectedIndex.value = index;
	} else {
		mapStore.clearLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`
		);
		selectedIndex.value = null;
	}
}
</script>

<template>
	<div v-if="activeChart === 'VennDiagram'" class="venndiagram">
		<div class="venndiagram-title">
			<h5>總機構數量</h5>
			<h6>{{ VennData.sum }} {{ chart_config.unit }}</h6>
		</div>
		<div class="venndiagram-chart">
			<div v-if="VennData.len === 2" class="circle">
				<svg
					viewBox="0 0 550 550"
					class="svg"
					xmlns="http://www.w3.org/2000/svg"
				>
					<g transform="translate(0 50)">
						<path
							data-name="0"
							fill="#D74F52"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '0' || selectedIndex === 0,
							}"
							d="M 275 145.09618943233
							A 150 150 0 1 1 275 404.90381056767
							A 150 150 0 0 0 275 145.09618943233
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(0)"
						/>
						<path
							data-name="1"
							fill="#764C7F"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '1' || selectedIndex === 1,
							}"
							d="M 275 145.09618943233
							A 150 150 0 1 0 275 404.90381056767
							A 150 150 0 0 1 275 145.09618943233
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(1)"
						/>
						<path
							data-name="2"
							fill="#483299"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '2' || selectedIndex === 2,
							}"
							d="M 275 145.09618943233
							A 150 150 0 0 0 275 404.90381056767
							A 150 150 0 0 0 275 145.09618943233
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(2)"
						/>
					</g>
				</svg>
			</div>
			<div v-if="VennData.len === 3" class="circle">
				<svg
					viewBox="0 0 550 550"
					class="svg"
					xmlns="http://www.w3.org/2000/svg"
					@mousemove="updateMouseLocation"
					@mouseleave="toggleActiveToNull"
				>
					<g transform="translate(0 50)">
						<path
							data-name="0"
							fill="#80E3D4"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '0' || selectedIndex === 0,
							}"
							d="M 125.6327126328 213.76275643042
							A 150 150 0 0 1 275 227.52551286084
							A 150 150 0 0 1 424.3672873672 213.76275643042
							A 150 150 0 1 0 125.6327126328 213.76275643042
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(0)"
						/>
						<path
							data-name="1"
							fill="#8CAE65"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '1' || selectedIndex === 1,
							}"
							d="M 125.6327126328 213.76275643042
							A 150 150 0 0 0 212.23525301124 336.23724356958
							A 150 150 0 0 0 275 472.47448713916
							A 150 150 0 1 1 125.6327126328 213.76275643042
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(1)"
						/>
						<path
							data-name="2"
							fill="#D6B059"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '2' || selectedIndex === 2,
							}"
							d="M 424.3672873672 213.76275643042
							A 150 150 0 0 1 337.76474698876 336.23724356958
							A 150 150 0 0 1 275 472.47448713916
							A 150 150 0 1 0 424.3672873672 213.76275643042
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(2)"
						/>

						<path
							data-name="3"
							fill="#D6B059"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '3' || selectedIndex === 3,
							}"
							d="M 125.6327126328 213.76275643042
							A 150 150 0 0 1 275 227.52551286084
							A 150 150 0 0 0 212.23525301124 336.23724356958
							A 150 150 0 0 1 125.6327126328 213.76275643042
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(3)"
						/>
						<path
							data-name="4"
							fill="#8CAE65"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '4' || selectedIndex === 4,
							}"
							d="M 275 227.52551286084
							A 150 150 0 0 1 424.3672873672 213.76275643042
							A 150 150 0 0 1 337.76474698876 336.23724356958
							A 150 150 0 0 0 275 227.52551286084
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(4)"
						/>
						<path
							data-name="5"
							fill="#80E3D4"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '5' || selectedIndex === 5,
							}"
							d="M 212.23525301124 336.23724356958
							A 150 150 0 0 0 337.76474698876 336.23724356958
							A 150 150 0 0 1 275 472.47448713916
							A 150 150 0 0 1 212.23525301124 336.23724356958
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(5)"
						/>
						<path
							data-name="6"
							fill="#eee"
							stroke="none"
							:class="{
								'active-block':
									targetArea === '6' || selectedIndex === 6,
							}"
							d="M 275 227.52551286084
							A 150 150 0 0 0 212.23525301124 336.23724356958
							A 150 150 0 0 0 337.76474698876 336.23724356958
							A 150 150 0 0 0 275 227.52551286084
							z
							"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
							@click="handleDataSelection(6)"
						/>
					</g>
				</svg>
			</div>
			<Teleport to="body">
				<!-- The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css -->
				<div
					v-if="targetArea"
					class="venndiagram-chart-info chart-tooltip"
					:style="tooltipPosition"
				>
					<h6>數量</h6>
					<span
						>{{ VennData.count[targetArea] }}
						{{ chart_config.unit }}</span
					>
				</div>
			</Teleport>
		</div>
	</div>
</template>

<style scoped lang="scss">
.venndiagram {
	max-height: 100%;
	position: relative;
	overflow-y: scroll;

	&-title {
		display: flex;
		flex-direction: column;
		justify-content: center;
		position: absolute;
		left: 0;
		top: 0;
		margin: 0.5rem 0 -0.5rem;

		h5 {
			color: var(--color-complement-text);
		}

		h6 {
			color: var(--color-complement-text);
			font-size: var(--font-m);
			font-weight: 400;
		}

		&-legend {
			display: flex;
			justify-content: space-between;

			div {
				position: relative;
				width: 3rem;
				margin: 0 4px;
				border-radius: 5px;
			}

			div:before {
				content: "";
				width: 3rem;
				height: var(--font-l);
				position: absolute;
				top: 0;
				left: 0;
				background: linear-gradient(
					270deg,
					rgba(40, 42, 44, 1),
					rgba(40, 42, 44, 0.2)
				);
			}

			p {
				color: var(--color-complement-text);
			}
		}
	}

	&-chart {
		display: flex;
		justify-content: center;

		svg {
			width: 80%;
			-webkit-transform: translateX(10%);
			-ms-transform: translateX(10%);
			transform: translateX(10%);
		}

		&-info {
			position: fixed;
			z-index: 20;
		}
	}
}
@keyframes ease-in {
	0% {
		transform: translateY(0px);
	}

	100% {
		transform: translateY(-7px);
	}
}
.active-block {
	animation-name: ease-in;
	animation-duration: 0.1s;
	animation-delay: 0.05s;
	animation-timing-function: linear;
	animation-fill-mode: forwards;
}
</style>
