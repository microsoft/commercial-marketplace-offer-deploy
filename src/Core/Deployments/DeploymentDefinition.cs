// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Text.Json.Serialization;
using Modm.Packaging;
using Modm.Serialization;

namespace Modm.Deployments
{
    /// <summary>
    /// Represents the contents of the manifest file contained in the archive file, e.g. the installer.zip inside the app.zip
    /// </summary>
	public class DeploymentDefinition
	{
        /// <summary>
        /// The fully qualified working directory that the deployment files are all located
        /// </summary>
        public string WorkingDirectory { get; set; }

        /// <summary>
        /// The relative (to the working directory) path to the main template
        /// </summary>
        public string MainTemplatePath { get; set; }

        /// <summary>
        /// The relative (to the working directory) path to the parameters file
        /// </summary>
        public string ParametersFilePath { get; set; }

        public string DeploymentType { get; set; }

        /// <summary>
        /// The source of the deployment definition
        /// </summary>
        public PackageUri Source { get; set; }

        /// <summary>
        /// The origional hash of the installer package content
        /// </summary>
        public string InstallerPackageHash { get; set; }

        [JsonConverter(typeof(DictionaryStringObjectJsonConverter))]
        public Dictionary<string, object> Parameters { get; set; }

        /// <summary>
        /// Gets the fully qualified directory path where the main template is located
        /// </summary>
        /// <returns></returns>
        public string GetMainTemplateDirectoryName()
        {
            return Path.GetDirectoryName(Path.Combine(WorkingDirectory, MainTemplatePath));
        }
    }
}

