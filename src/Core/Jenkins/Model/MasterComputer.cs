// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Text.Json.Serialization;

namespace Modm.Jenkins.Model
{
    /// <summary>
    /// Generated from JSON using http://localhost:8080/computer/(built-in)/api/json
    /// </summary>
    public class MasterComputer
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("actions")]
        public List<Action> Actions { get; set; }

        [JsonPropertyName("assignedLabels")]
        public List<AssignedLabel> AssignedLabels { get; set; }

        [JsonPropertyName("description")]
        public string Description { get; set; }

        [JsonPropertyName("displayName")]
        public string DisplayName { get; set; }

        [JsonPropertyName("executors")]
        public List<Executor> Executors { get; set; }

        [JsonPropertyName("icon")]
        public string Icon { get; set; }

        [JsonPropertyName("iconClassName")]
        public string IconClassName { get; set; }

        [JsonPropertyName("idle")]
        public bool Idle { get; set; }

        [JsonPropertyName("jnlpAgent")]
        public bool JnlpAgent { get; set; }

        [JsonPropertyName("launchSupported")]
        public bool LaunchSupported { get; set; }

        [JsonPropertyName("loadStatistics")]
        public LoadStatistics LoadStatistics { get; set; }

        [JsonPropertyName("manualLaunchAllowed")]
        public bool ManualLaunchAllowed { get; set; }

        [JsonPropertyName("monitorData")]
        public MonitorData MonitorData { get; set; }

        [JsonPropertyName("numExecutors")]
        public int NumExecutors { get; set; }

        [JsonPropertyName("offline")]
        public bool Offline { get; set; }

        [JsonPropertyName("offlineCause")]
        public object OfflineCause { get; set; }

        [JsonPropertyName("offlineCauseReason")]
        public string OfflineCauseReason { get; set; }

        [JsonPropertyName("oneOffExecutors")]
        public List<object> OneOffExecutors { get; set; }

        [JsonPropertyName("temporarilyOffline")]
        public bool TemporarilyOffline { get; set; }
    }

    public class Action
    {
    }

    public class Executor
    {
    }

    public class HudsonNodeMonitorsClockMonitor
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("diff")]
        public int Diff { get; set; }
    }

    public class HudsonNodeMonitorsDiskSpaceMonitor
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("timestamp")]
        public long Timestamp { get; set; }

        [JsonPropertyName("path")]
        public string Path { get; set; }

        [JsonPropertyName("size")]
        public long Size { get; set; }
    }

    public class HudsonNodeMonitorsResponseTimeMonitor
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("timestamp")]
        public long Timestamp { get; set; }

        [JsonPropertyName("average")]
        public int Average { get; set; }
    }

    public class HudsonNodeMonitorsSwapSpaceMonitor
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("availablePhysicalMemory")]
        public long AvailablePhysicalMemory { get; set; }

        [JsonPropertyName("availableSwapSpace")]
        public int AvailableSwapSpace { get; set; }

        [JsonPropertyName("totalPhysicalMemory")]
        public long TotalPhysicalMemory { get; set; }

        [JsonPropertyName("totalSwapSpace")]
        public int TotalSwapSpace { get; set; }
    }

    public class HudsonNodeMonitorsTemporarySpaceMonitor
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }

        [JsonPropertyName("timestamp")]
        public long Timestamp { get; set; }

        [JsonPropertyName("path")]
        public string Path { get; set; }

        [JsonPropertyName("size")]
        public long Size { get; set; }
    }

    public class LoadStatistics
    {
        [JsonPropertyName("_class")]
        public string Class { get; set; }
    }

    public class MonitorData
    {
        [JsonPropertyName("hudson.node_monitors.SwapSpaceMonitor")]
        public HudsonNodeMonitorsSwapSpaceMonitor HudsonNodeMonitorsSwapSpaceMonitor { get; set; }

        [JsonPropertyName("hudson.node_monitors.TemporarySpaceMonitor")]
        public HudsonNodeMonitorsTemporarySpaceMonitor HudsonNodeMonitorsTemporarySpaceMonitor { get; set; }

        [JsonPropertyName("hudson.node_monitors.DiskSpaceMonitor")]
        public HudsonNodeMonitorsDiskSpaceMonitor HudsonNodeMonitorsDiskSpaceMonitor { get; set; }

        [JsonPropertyName("hudson.node_monitors.ArchitectureMonitor")]
        public string HudsonNodeMonitorsArchitectureMonitor { get; set; }

        [JsonPropertyName("hudson.node_monitors.ResponseTimeMonitor")]
        public HudsonNodeMonitorsResponseTimeMonitor HudsonNodeMonitorsResponseTimeMonitor { get; set; }

        [JsonPropertyName("hudson.node_monitors.ClockMonitor")]
        public HudsonNodeMonitorsClockMonitor HudsonNodeMonitorsClockMonitor { get; set; }
    }

}

