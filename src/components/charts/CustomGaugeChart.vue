<script setup>
import { ref } from "vue";
const props = defineProps([
	"chart_config",
	"activeChart",
	"series",
	"map_config",
]);
const R = 180;
const Rx = 275;
const Ry = 240;
const Llength = 15;
const mousePosition = ref({ x: null, y: null });
let targetvalue = ref(0);
let selectedvalue = ref(0);

function toggleActive(e) {
	targetvalue.value = e.target.dataset.name;
}
function toggleActiveToNull() {
	// targetvalue.value = 0;
}
function updateMouseLocation(e) {
	mousePosition.value.x = e.pageX;
	mousePosition.value.y = e.pageY;
}
function calGradComp(chart_config) {
	const percent = chart_config["percent"];
	const color = chart_config["grad_color"];
	const res = [];
	res.push([color[0], "0%"]);
	for (let i = 0; i < percent.length - 1; i++) {
		res.push([
			color[i],
			(
				(parseFloat(percent[i]) + parseFloat(percent[i + 1])) /
				2
			).toString() + "%",
		]);
	}
	return res;
}
function calAngle(percentvalue) {
	return percentvalue * 1.8;
}
function calPt(angle, radius) {
	const rad = (angle * Math.PI) / 180;
	return { x: Rx - radius * Math.cos(rad), y: Ry - radius * Math.sin(rad) };
}
function calLine(angle, length) {
	const rad = (angle * Math.PI) / 180;
	return { x: length * Math.cos(rad), y: length * Math.sin(rad) };
}
function calTri(value, chart_config) {
	const standards = chart_config["standards"];
	const percents = chart_config["percent"];

	let percentvalue = 0;
	if (value <= standards[0]) percentvalue = 0;
	else if (value >= standards[standards.length - 1]) percentvalue = 180;
	else {
		let idx = 0;
		while (standards[idx] < value) idx += 1;
		percentvalue =
			((value - standards[idx - 1]) /
				(standards[idx] - standards[idx - 1])) *
				(percents[idx] - percents[idx - 1]) +
			percents[idx - 1];
	}
	// console.log(percentvalue, value);
	let triAngle = calAngle(percentvalue);
	let triPt = calPt(triAngle, R - Llength);
	let res = "";
	let topAngle = 20;
	let sideLength = 30;
	res += "M " + triPt.x + " " + triPt.y + "\n";
	res +=
		"l " +
		calLine(triAngle + topAngle, sideLength).x +
		" " +
		calLine(triAngle + topAngle, sideLength).y +
		"\n";
	res +=
		"l " +
		calLine(
			triAngle - 90,
			2 * sideLength * Math.sin((topAngle * Math.PI) / 180)
		).x +
		" " +
		calLine(
			triAngle - 90,
			2 * sideLength * Math.sin((topAngle * Math.PI) / 180)
		).y +
		"\n";
	// res += "z\n";
	return res;
}
function CheckpathD(startpercent, endpercent, chart_config) {
	let startAngle = calAngle(startpercent, chart_config);
	let endAngle = calAngle(endpercent, chart_config);

	let startPt = calPt(startAngle, R);
	let endPt = calPt(endAngle, R);
	let startbotPt = calPt(startAngle, R - Llength);
	let line = calLine(endAngle, Llength);

	let res = "";
	res += "M " + startPt.x + " " + startPt.y + "\n";
	res += "A " + R + " " + R + " 0 0 1 " + endPt.x + " " + endPt.y + "\n";
	res += "l " + line.x + " " + line.y + "\n";
	res +=
		"A " +
		(R - Llength) +
		" " +
		(R - Llength) +
		" 0 0 0 " +
		startbotPt.x +
		" " +
		startbotPt.y +
		"\n";
	res += "z\n";
	return res;
}
function handleTable(list) {
	let res = [[], []];
	// 2 by N/2
	for (let i = 0; i <= list.length / 2; i++) {
		res[0].push([2 * i, list[2 * i]]);
		if (2 * i + 1 < list.length) res[1].push([2 * i + 1, list[2 * i + 1]]);
	}
	return res;
}
</script>
<template>
	<div v-if="activeChart === 'CustomGaugeChart'" class="air">
		<g>
			<svg
				viewBox="0 0 550 260"
				xmlns="http://www.w3.org/2000/svg"
				class="initial-animation-1"
			>
				<defs>
					<linearGradient
						:id="'grad1-' + chart_config.name"
						x1="0%"
						y1="0%"
						x2="100%"
						y2="0%"
					>
						<stop
							v-for="grad in calGradComp(chart_config)"
							:offset="grad[1]"
							:style="{
								'stop-color': grad[0],
								'stop-opacity': 1,
							}"
						/>
					</linearGradient>
				</defs>
				<path
					:fill="'url(#grad1-' + chart_config.name + ')'"
					:d="CheckpathD(0, 100, chart_config)"
				/>
				<path
					fill="#ddd"
					stroke="#ddd"
					:d="calTri(series[0].data[targetvalue], chart_config)"
				/>
				<text
					:x="Rx"
					:y="Ry"
					dominant-baseline="start"
					text-anchor="middle"
					class="small"
				>
					<tspan>
						{{
							chart_config.name +
							" (" +
							chart_config.categories[targetvalue] +
							")"
						}}
					</tspan>
				</text>
				<text
					:x="Rx - 10"
					:y="Ry - 35"
					dominant-baseline="start"
					text-anchor="middle"
					class="heavy"
				>
					{{ series[0].data[targetvalue] }}
				</text>
				<text
					:x="Rx + 90"
					:y="Ry - 30"
					dominant-baseline="start"
					text-anchor="middle"
					class="small"
				>
					<tspan>{{ chart_config.unit }}</tspan>
				</text>
			</svg>
			<div class="tablediv">
				<table>
					<tbody>
						<tr
							v-for="subTable in handleTable(
								chart_config.categories
							)"
						>
							<td
								v-for="submatter in subTable"
								:data-name="submatter[0].toString()"
								@mouseenter="toggleActive"
								@mouseleave="toggleActiveToNull"
								:class="'initial-animation-' + submatter[0]"
							>
								<div class="mattertable">
									{{ submatter[1] }}
								</div>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</g>
	</div>
</template>
<style scoped lang="scss">
.air {
	user-select: none;
}
.small {
	font: 16px sans-serif;
	fill: #ddd;
}
.smallsub {
	font: 12px sans-serif;
	fill: #ddd;
}
.heavy {
	font: 60px sans-serif;
	fill: #ddd;
}
.tablediv {
	display: flex;
	justify-content: center;
}
table {
	/* border-collapse: collapse; */
	/* width: 90%; */
	border-spacing: 3px;
	width: 90%;
	table-layout: fixed;
}
th,
td {
	width: 30%;
	/* border: 2px solid #666; */
	padding: 6px;
	text-align: left;
	border-radius: 5px;
	font-size: 12px;
	background-color: #444444;
}
td:hover {
	background-color: #111111;
	cursor: default;
}
.mattertable {
	width: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
}
.mattervalue {
	display: flex;
	justify-content: space-around;
	align-items: end;
	font-size: 10px;
}
.mattername {
	font-size: 14px;
	font-family: sans-serif;
}
@keyframes ease-in {
	0% {
		opacity: 0;
	}

	100% {
		opacity: 1;
	}
}

@for $i from 1 through 20 {
	.initial-animation-#{$i} {
		animation-name: ease-in;
		animation-duration: 0.3s;
		animation-delay: 0.1s * ($i - 1);
		animation-timing-function: linear;
		animation-fill-mode: forwards;
		opacity: 0;
	}
}
</style>
