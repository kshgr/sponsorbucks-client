import React from "react";
import {
  AbsoluteFill,
  Easing,
  interpolate,
  spring,
  useCurrentFrame,
  useVideoConfig,
} from "remotion";

export type BarDatum = {
  label: string;
  value: number;
  color: string;
};

export type AnimatedBarChartProps = {
  data: BarDatum[];
  title?: string;
  subtitle?: string;
};

const CARD_BG = "rgba(10, 14, 26, 0.88)";
const GRID_COLOR = "rgba(255, 255, 255, 0.08)";
const AXIS_COLOR = "rgba(255, 255, 255, 0.26)";

export const AnimatedBarChart: React.FC<AnimatedBarChartProps> = ({
  data,
  title = "Performance Snapshot",
  subtitle = "Five bars animating in sequence",
}) => {
  const frame = useCurrentFrame();
  const { fps } = useVideoConfig();

  const cardOpacity = interpolate(frame, [0, 18], [0, 1], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
    easing: Easing.out(Easing.quad),
  });
  const cardY = interpolate(frame, [0, 18], [24, 0], {
    extrapolateLeft: "clamp",
    extrapolateRight: "clamp",
    easing: Easing.out(Easing.quad),
  });

  return (
    <AbsoluteFill
      style={{
        background:
          "radial-gradient(circle at top, #1a2442 0%, #0a1020 48%, #050814 100%)",
        color: "white",
        fontFamily:
          '"Inter", "Avenir Next", "Segoe UI", "Helvetica Neue", Arial, sans-serif',
        padding: 56,
      }}
    >
      <div
        style={{
          position: "absolute",
          inset: 0,
          background:
            "linear-gradient(135deg, rgba(76, 99, 255, 0.14), transparent 42%, rgba(26, 231, 198, 0.08) 84%)",
        }}
      />

      <div
        style={{
          position: "relative",
          zIndex: 1,
          display: "flex",
          flexDirection: "column",
          gap: 28,
          width: "100%",
          height: "100%",
        }}
      >
        <div style={{ maxWidth: 720 }}>
          <div
            style={{
              letterSpacing: "0.24em",
              textTransform: "uppercase",
              color: "rgba(255,255,255,0.64)",
              fontSize: 18,
              marginBottom: 10,
            }}
          >
            Animated chart
          </div>
          <h1
            style={{
              fontSize: 72,
              lineHeight: 1.02,
              margin: 0,
              fontWeight: 800,
            }}
          >
            {title}
          </h1>
          <p
            style={{
              margin: "16px 0 0",
              fontSize: 28,
              lineHeight: 1.35,
              color: "rgba(255,255,255,0.74)",
              maxWidth: 680,
            }}
          >
            {subtitle}
          </p>
        </div>

        <div
          style={{
            flex: 1,
            minHeight: 0,
            borderRadius: 36,
            border: "1px solid rgba(255,255,255,0.12)",
            background: CARD_BG,
            boxShadow: "0 40px 120px rgba(0, 0, 0, 0.45)",
            backdropFilter: "blur(18px)",
            padding: 36,
            opacity: cardOpacity,
            transform: `translateY(${cardY}px)`,
          }}
        >
          <div
            style={{
              display: "flex",
              flexDirection: "column",
              height: "100%",
              gap: 22,
            }}
          >
            <div
              style={{
                display: "flex",
                alignItems: "baseline",
                justifyContent: "space-between",
                gap: 20,
              }}
            >
              <div>
                <div
                  style={{
                    color: "rgba(255,255,255,0.62)",
                    fontSize: 18,
                    marginBottom: 8,
                  }}
                >
                  Monthly growth
                </div>
                <div style={{ fontSize: 36, fontWeight: 700 }}>
                  5-bar animated comparison
                </div>
              </div>

              <div
                style={{
                  display: "flex",
                  alignItems: "center",
                  gap: 10,
                  color: "rgba(255,255,255,0.7)",
                  fontSize: 18,
                }}
              >
                <span
                  style={{
                    width: 12,
                    height: 12,
                    borderRadius: 999,
                    background: "#69f0ae",
                    boxShadow: "0 0 18px rgba(105, 240, 174, 0.85)",
                  }}
                />
                Live motion
              </div>
            </div>

            <div
              style={{
                flex: 1,
                minHeight: 0,
                display: "grid",
                gridTemplateColumns: "88px 1fr",
                gap: 18,
              }}
            >
              <div
                style={{
                  position: "relative",
                  color: "rgba(255,255,255,0.4)",
                  fontSize: 16,
                }}
              >
                {[100, 75, 50, 25, 0].map((tick, index) => (
                  <div
                    key={tick}
                    style={{
                      position: "absolute",
                      top: `${index * 25}%`,
                      right: 12,
                      transform: "translateY(-50%)",
                    }}
                  >
                    {tick}
                  </div>
                ))}
              </div>

              <div style={{ position: "relative" }}>
                {[0, 1, 2, 3, 4].map((index) => (
                  <div
                    key={index}
                    style={{
                      position: "absolute",
                      inset: 0,
                      backgroundImage: `linear-gradient(to top, ${GRID_COLOR} 1px, transparent 1px)`,
                      backgroundSize: "100% 25%",
                      opacity: 0.8,
                    }}
                  />
                ))}

                <div
                  style={{
                    position: "relative",
                    zIndex: 1,
                    display: "grid",
                    gridTemplateColumns: `repeat(${data.length}, minmax(0, 1fr))`,
                    alignItems: "end",
                    height: "100%",
                    gap: 20,
                    padding: "0 10px 32px 10px",
                  }}
                >
                  {data.map((bar, index) => {
                    const progress = spring({
                      frame,
                      fps,
                      delay: index * 6,
                      config: {
                        damping: 16,
                        stiffness: 120,
                        mass: 0.9,
                      },
                    });
                    const barScale = interpolate(progress, [0, 1], [0.05, 1], {
                      extrapolateLeft: "clamp",
                      extrapolateRight: "clamp",
                    });
                    const labelOpacity = interpolate(
                      frame,
                      [18 + index * 6, 28 + index * 6],
                      [0, 1],
                      {
                        extrapolateLeft: "clamp",
                        extrapolateRight: "clamp",
                      }
                    );

                    return (
                      <div
                        key={bar.label}
                        style={{
                          display: "flex",
                          flexDirection: "column",
                          justifyContent: "end",
                          gap: 16,
                          height: "100%",
                        }}
                      >
                        <div
                          style={{
                            display: "flex",
                            alignItems: "flex-end",
                            justifyContent: "center",
                            height: "100%",
                            minHeight: 320,
                          }}
                        >
                          <div
                            style={{
                              width: "100%",
                              maxWidth: 132,
                              height: `${bar.value}%`,
                              transform: `scaleY(${barScale})`,
                              transformOrigin: "bottom",
                              borderRadius: "20px 20px 10px 10px",
                              background: `linear-gradient(180deg, ${bar.color} 0%, rgba(255,255,255,0.12) 160%)`,
                              boxShadow: `0 18px 40px ${bar.color}44`,
                              border: "1px solid rgba(255,255,255,0.1)",
                              position: "relative",
                              overflow: "hidden",
                            }}
                          >
                            <div
                              style={{
                                position: "absolute",
                                inset: 0,
                                background:
                                  "linear-gradient(135deg, rgba(255,255,255,0.22), transparent 42%, transparent 58%, rgba(255,255,255,0.14))",
                                opacity: 0.8,
                              }}
                            />
                            <div
                              style={{
                                position: "absolute",
                                top: 14,
                                left: "50%",
                                transform: "translateX(-50%)",
                                fontSize: 20,
                                fontWeight: 700,
                                opacity: labelOpacity,
                                textShadow: "0 2px 8px rgba(0,0,0,0.45)",
                              }}
                            >
                              {bar.value.toFixed(0)}%
                            </div>
                          </div>
                        </div>

                        <div
                          style={{
                            textAlign: "center",
                            fontSize: 20,
                            fontWeight: 600,
                            color: "rgba(255,255,255,0.78)",
                            letterSpacing: "0.02em",
                            opacity: labelOpacity,
                          }}
                        >
                          {bar.label}
                        </div>
                      </div>
                    );
                  })}
                </div>

                <div
                  style={{
                    position: "absolute",
                    left: 0,
                    right: 0,
                    bottom: 0,
                    height: 1,
                    background: AXIS_COLOR,
                  }}
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </AbsoluteFill>
  );
};
