import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import fs from 'fs';
import path from 'path';
import dns from 'dns';

const aspnetcoreHttps = function() {
  // this is going to use the certificates created by asp.net core for local development
  const baseFolder =
    process.env.APPDATA !== undefined && process.env.APPDATA !== ''
      ? `${process.env.APPDATA}/ASP.NET/https`
      : `${process.env.HOME}/.aspnet/https`;

  const certificateArg = process.argv.map(arg => arg.match(/--name=(?<value>.+)/i)).filter(Boolean)[0];
  const certificateName = certificateArg ? certificateArg.groups.value : process.env.npm_package_name;

  if (!certificateName) {
    console.error('Invalid certificate name. Run this script in the context of an npm/yarn script or pass --name=<<app>> explicitly.')
    process.exit(-1);
  }

  const certFilePath = path.join(baseFolder, `${certificateName}.pem`);
  const keyFilePath = path.join(baseFolder, `${certificateName}.key`);

  return {
    key: fs.readFileSync(keyFilePath),
    cert: fs.readFileSync(certFilePath)
  }
}

dns.setDefaultResultOrder('verbatim')

export default defineConfig(({command, mode}) => {
  
  let config = {
    server: {
      host: 'localhost',
      port: 44482,
    },
    build: {
      outDir: 'build',
    },
    plugins: [react()],
  };

  if (mode === 'development') {
    config.server.https = aspnetcoreHttps()
  }

  return config;
});