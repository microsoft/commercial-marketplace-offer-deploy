// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
﻿using System;
using System.Text.Json;
using Modm.Serialization;

namespace Modm.Deployments
{
    /// <summary>
    /// Writes out parameters for a deployment as tfvars file
    /// </summary>
	class TerraformParametersFile : IDeploymentParametersFile
	{
        public string FullPath
        {
            get {  return Path.Combine(destinationDirectory, FileName); }
        }

        /// <summary>
        /// This file name will automatically be picked up by terraform <see cref="https://developer.hashicorp.com/terraform/language/values/variables"/>
        /// </summary>
        public const string FileName = "terraform.tfvars.json";
        private readonly string destinationDirectory;

        public TerraformParametersFile(string destinationDirectory)
		{
            this.destinationDirectory = destinationDirectory;
        }

        public async Task Write(IDictionary<string, object> parameters)
        {
            var json = JsonSerializer.Serialize(parameters, new JsonSerializerOptions
            {
                Converters = { new DictionaryStringObjectJsonConverter() }
            });
            await File.WriteAllTextAsync(FullPath, json);
        }
    }
}

