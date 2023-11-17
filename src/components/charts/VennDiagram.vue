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
// const selectedIndex = ref(null);

// Parse Venn Diagram Data
const VennData = computed(() => {
	let readData = {};
	let count = {};
	const sum = props.series.length;

	let len = props.series[0].data.length;
	if (len === 2) {
		for (let i = 0; i < 3; i++) {
			count[i.toString()] = 0;
		}
	} else if (len === 3) {
		for (let i = 0; i < 7; i++) {
			count[i.toString()] = 0;
		}
	}

	for (let i = 0; i < props.series.length; i++) {
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

function toggleActive(len) {
	const mouseX = refMousePosition.value.x;
	const mouseY = refMousePosition.value.y;
	if (len === 2) {
		const c1 = [200, 275];
		const c2 = [350, 275];
		const d1 = ((mouseX - c1[0]) ** 2 + (mouseY - c1[1]) ** 2) ** 0.5;
		const d2 = ((mouseX - c2[0]) ** 2 + (mouseY - c2[1]) ** 2) ** 0.5;
		if ((d1 <= 150) & (d2 <= 150)) {
			targetArea.value = "2";
		} else if ((d1 <= 150) & (d2 > 150)) {
			targetArea.value = "0";
		} else if ((d1 > 150) & (d2 <= 150)) {
			targetArea.value = "1";
		}
	} else if (len === 3) {
		const c1 = [275, 200];
		const c2 = [188.3975, 350];
		const c3 = [361.6025, 350];
		const d1 = ((mouseX - c1[0]) ** 2 + (mouseY - c1[1]) ** 2) ** 0.5;
		const d2 = ((mouseX - c2[0]) ** 2 + (mouseY - c2[1]) ** 2) ** 0.5;
		const d3 = ((mouseX - c3[0]) ** 2 + (mouseY - c3[1]) ** 2) ** 0.5;
		if ((d1 <= 150) & (d2 <= 150) & (d3 <= 150)) {
			targetArea.value = "6";
		} else if ((d1 <= 150) & (d2 > 150) & (d3 > 150)) {
			targetArea.value = "0";
		} else if ((d1 > 150) & (d2 <= 150) & (d3 > 150)) {
			targetArea.value = "1";
		} else if ((d1 > 150) & (d2 > 150) & (d3 <= 150)) {
			targetArea.value = "2";
		} else if ((d1 <= 150) & (d2 <= 150) & (d3 > 150)) {
			targetArea.value = "3";
		} else if ((d1 <= 150) & (d2 > 150) & (d3 <= 150)) {
			targetArea.value = "4";
		} else if ((d1 > 150) & (d2 <= 150) & (d3 <= 150)) {
			targetArea.value = "5";
		} else {
			targetArea.value = null;
		}
	}
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

/*
function handleDataSelection(index) {
	if (!props.chart_config.map_filter) {
		return;
	}
	if (index !== selectedIndex.value) {
		mapStore.addLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`,
			props.chart_config.map_filter[0],
			props.chart_config.map_filter[1][index]
		);
		selectedIndex.value = index;
	} else {
		mapStore.clearLayerFilter(
			`${props.map_config[0].index}-${props.map_config[0].type}`
		);
		selectedIndex.value = null;
	}
}
*/
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
					<g transform="translate(350 275)">
						<circle
							cx="0"
							cy="0"
							r="150"
							fill="#D74F52"
							opacity="0.6"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
						/>
					</g>
					<g transform="translate(200 275)">
						<circle
							cx="0"
							cy="0"
							r="150"
							fill="#764C7F"
							opacity="0.6"
							@mouseenter="toggleActive"
							@mousemove="updateMouseLocation"
							@mouseleave="toggleActiveToNull"
						/>
					</g>
				</svg>
			</div>
			<div v-if="VennData.len === 3" class="circle">
				<svg
					viewBox="0 0 550 550"
					class="svg"
					xmlns="http://www.w3.org/2000/svg"
					@mousemove="
						updateMouseLocation($event);
						toggleActive(3);
					"
					@mouseleave="toggleActiveToNull"
				>
					<g transform="translate(275 200)">
						<circle
							cx="0"
							cy="0"
							r="150"
							fill="#80E3D4"
							opacity="0.6"
						/>
					</g>
					<g transform="translate(188.3975 350)">
						<circle
							cx="0"
							cy="0"
							r="150"
							fill="#8CAE65"
							opacity="0.6"
						/>
					</g>
					<g transform="translate(361.6025 350)">
						<circle
							cx="0"
							cy="0"
							r="150"
							fill="#D6B059"
							opacity="0.6"
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
</style>
