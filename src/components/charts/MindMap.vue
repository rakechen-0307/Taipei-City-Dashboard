<script setup>
import { ref, computed } from "vue";
import { useMapStore } from "../../store/mapStore";

const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);
const mapStore = useMapStore();

const targetCategory = ref(null);
// const AreaColor = ref(props.chart_config.color[0]);
const mousePosition = ref({ x: null, y: null });
const selectedIndex = ref(null);

// Parse Mind Map Data
const MindData = computed(() => {
	let sum = 0;
	let highest = 0;
	let typeCount = {};
	const count = props.series[0].data.length;
	let config = [];

	for (let i = 0; i < count; i++) {
		typeCount[props.series[0].data[i].x] = props.series[0].data[i].y;
		sum = sum + props.series[0].data[i].y;
		if (props.series[0].data[i].y > highest) {
			highest = props.series[0].data[i].y;
		}
	}

	const rad = (2 * Math.PI) / count;
	let x1, x2, y1, y2, r, cx, cy, tx, ty, color, name;
	for (let i = 0; i < count; i++) {
		x1 = 90 * Math.cos(rad * i);
		y1 = 90 * Math.sin(rad * i);
		x2 = 150 * Math.cos(rad * i);
		y2 = 150 * Math.sin(rad * i);
		r = 35 + 10 * Math.log10(props.series[0].data[i].y);
		// eslint-disable-next-line no-console
		console.log(r);
		cx = (150 + r) * Math.cos(rad * i);
		cy = (150 + r) * Math.sin(rad * i);
		tx = 300 * Math.cos(rad * i);
		ty = 300 * Math.sin(rad * i);
		color = props.chart_config.color[i];
		name = props.chart_config.map_filter[1][i];
		config.push([x1, y1, x2, y2, r, cx, cy, tx, ty, color, name, i]);
	}

	const output = {
		sum: sum,
		count: count,
		typeCount: typeCount,
		config: config,
	};
	return output;
});

const tooltipPosition = computed(() => {
	return {
		left: `${mousePosition.value.x - 10}px`,
		top: `${mousePosition.value.y - 54}px`,
	};
});

function toggleActive(e) {
	targetCategory.value = e.target.dataset.name;
}
function toggleActiveToNull() {
	targetCategory.value = null;
}
function updateMouseLocation(e) {
	mousePosition.value.x = e.pageX;
	mousePosition.value.y = e.pageY;
}

function handleDataSelection(index) {
	// eslint-disable-next-line no-console
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
</script>

<template>
	<div v-if="activeChart === 'MindMap'" class="mindmap">
		<div class="mindmap-chart">
			<svg
				viewBox="0 0 650 600"
				class="svg"
				xmlns="http://www.w3.org/2000/svg"
			>
				<g transform="translate(250 335)">
					<circle
						cx="0"
						cy="0"
						r="90"
						stroke="#767575"
						stroke-width="5"
						fill-opacity="0"
					/>
					<text
						x="0"
						y="-20"
						text-anchor="middle"
						font-size="24pt"
						fill="#767575"
						alignment-baseline="middle"
					>
						總數量
					</text>
					<text
						x="0"
						y="20"
						text-anchor="middle"
						font-size="24pt"
						fill="#767575"
						alignment-baseline="middle"
					>
						{{ MindData.sum }} {{ chart_config.unit }}
					</text>
					<line
						:x1="p[0]"
						:y1="p[1]"
						:x2="p[2]"
						:y2="p[3]"
						style="stroke: #767575; stroke-width: 2"
						:class="'animation-' + (2 * p[11] + 1)"
						v-for="p in MindData.config"
						:key="p"
					></line>
					<circle
						:data-name="p[10]"
						:cx="p[5]"
						:cy="p[6]"
						:r="p[4]"
						stroke="#767575"
						stroke-width="5"
						:fill="p[9]"
						:class="'animation-' + (2 * p[11] + 2)"
						@mouseenter="toggleActive"
						@mouseleave="toggleActiveToNull"
						@mousemove="updateMouseLocation"
						@click="handleDataSelection(p[11])"
						v-for="p in MindData.config"
						:key="p"
					></circle>
				</g>
			</svg>
			<Teleport to="body">
				<!-- The class "chart-tooltip" could be edited in /assets/styles/chartStyles.css -->
				<div
					v-if="targetCategory"
					class="mindmap-chart-info chart-tooltip"
					:style="tooltipPosition"
				>
					<h6>{{ targetCategory }}</h6>
					<span
						>{{ MindData.typeCount[targetCategory] }}
						{{ chart_config.unit }}</span
					>
				</div>
			</Teleport>
		</div>
	</div>
</template>

<style scoped lang="scss">
.mindmap {
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
			width: 75%;
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
		opacity: 0;
	}

	100% {
		opacity: 1;
	}
}
@for $i from 1 through 40 {
	.animation-#{$i} {
		animation-name: ease-in;
		animation-duration: 0.25s;
		animation-delay: 0.15s * ($i - 1);
		animation-timing-function: linear;
		animation-fill-mode: forwards;
		opacity: 0;
	}
}
</style>
