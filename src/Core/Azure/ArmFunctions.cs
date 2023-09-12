using System;
using Azure.Deployments.Expression.Expressions;
using Newtonsoft.Json.Linq;

namespace Modm.Azure
{
    static class ArmFunctions
    {
        public const string FunctionName_UniqueString = "uniqueString";

        public static string? UniqueString(params string[] values)
        {
            var parameters = values.Select(arg => new FunctionArgument(JToken.FromObject(arg))).ToArray();
            var result = ExpressionBuiltInFunctions.Functions.EvaluateFunction(FunctionName_UniqueString, parameters, null);
            return result.Value<string>();
        }
    }
}

