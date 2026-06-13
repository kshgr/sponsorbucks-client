import React from "react";
import { Composition } from "remotion";
import { AnimatedBarChart, AnimatedBarChartProps } from "./AnimatedBarChart";

const demoData: AnimatedBarChartProps["data"] = [
  { label: "Alpha", value: 42, color: "#5EEAD4" },
  { label: "Beta", value: 68, color: "#60A5FA" },
  { label: "Gamma", value: 55, color: "#A78BFA" },
  { label: "Delta", value: 84, color: "#F59E0B" },
  { label: "Epsilon", value: 73, color: "#F472B6" },
];

export const RemotionRoot: React.FC = () => {
  return (
    <Composition
      id="AnimatedBarChart"
      component={AnimatedBarChart}
      durationInFrames={180}
      fps={30}
      width={1280}
      height={720}
      defaultProps={
        {
          data: demoData,
          title: "Five Bars, One Story",
          subtitle: "A clean motion chart with staggered bar entrances and labels.",
        } satisfies AnimatedBarChartProps
      }
    />
  );
};
