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
const RadialHist = computed(() => {
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
	let dist,
		density,
		s1_x,
		s1_y,
		s2_x,
		s2_y,
		d,
		e1_x,
		e1_y,
		e2_x,
		e2_y,
		color,
		ponits_str,
		name;
	for (let i = 0; i < count; i++) {
		dist = 85;
		density = 2 / 8;
		s1_x = dist * Math.cos(rad * (i - density / 1.5));
		s1_y = dist * Math.sin(rad * (i - density / 1.5));
		s2_x = dist * Math.cos(rad * (i + density / 1.5));
		s2_y = dist * Math.sin(rad * (i + density / 1.5));
		// d = 20 + 100 * Math.log10(props.series[0].data[i].y);
		d = 1200 * (props.series[0].data[i].y / sum);
		e2_x = (dist + d) * Math.cos(rad * (i - density / 1.2));
		e2_y = (dist + d) * Math.sin(rad * (i - density / 1.2));
		e1_x = (dist + d) * Math.cos(rad * (i + density / 1.2));
		e1_y = (dist + d) * Math.sin(rad * (i + density / 1.2));
		// r = 35 + 10 * Math.log10(props.series[0].data[i].y);
		// // eslint-disable-next-line no-console
		// console.log(r);
		// cx = (150 + r) * Math.cos(rad * i);
		// cy = (150 + r) * Math.sin(rad * i);
		// tx = 300 * Math.cos(rad * i);
		// ty = 300 * Math.sin(rad * i);
		color = props.chart_config.color[i];
		name = props.chart_config.map_filter[1][i];
		// config.push([x1, y1, x2, y2, r, cx, cy, tx, ty, color, name, i]);
		ponits_str =
			s1_x.toString() +
			"," +
			s1_y.toString() +
			" " +
			s2_x.toString() +
			"," +
			s2_y.toString() +
			" " +
			e1_x.toString() +
			"," +
			e1_y.toString() +
			" " +
			e2_x.toString() +
			"," +
			e2_y.toString();
		config.push([ponits_str, name, color, i]);
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
	<div v-if="activeChart === 'RadialHistogram'" class="mindmap">
		<div class="mindmap-chart">
			<svg
				viewBox="0 0 650 600"
				class="svg"
				xmlns="http://www.w3.org/2000/svg"
			>
				<g transform="translate(250 335)">
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
						{{ RadialHist.sum }} {{ chart_config.unit }}
					</text>
					<polygon
						:data-name="p[1]"
						:points="p[0]"
						stroke="#767575"
						stroke-width="3"
						:fill="p[2]"
						:class="'animation-' + (2 * p[3] + 2)"
						@mouseenter="toggleActive"
						@mouseleave="toggleActiveToNull"
						@mousemove="updateMouseLocation"
						@click="handleDataSelection(p[3])"
						v-for="p in RadialHist.config"
						:key="p"
					></polygon>
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
						>{{ RadialHist.typeCount[targetCategory] }}
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
